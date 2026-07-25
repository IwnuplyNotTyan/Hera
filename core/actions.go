package generate

func (m *Model) doConfirm() *Model {
	if len(m.Players) == 0 || m.CurrentPlayer >= len(m.Players) {
		return m
	}
	p := Point{m.CursorX, m.CursorY}
	current := m.Players[m.CurrentPlayer]
	wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

	if m.UltMode && !m.Shot {
		m = m.doUlt()
		cur := m.Players[m.CurrentPlayer]
		m.CursorX = cur.X
		m.CursorY = cur.Y
	} else if m.PushStrikeMode && !m.Shot {
		m = m.doPushStrike()
		cur := m.Players[m.CurrentPlayer]
		m.CursorX = cur.X
		m.CursorY = cur.Y
	} else if m.ShootMode && !m.Shot {
		if HasEffect(m.Players[m.CurrentPlayer].Effects, EffectSmoke) {
			m.Shot = true
			m.ShootMode = false
			return m
		}
		if !m.IsInRange(m.CursorX, m.CursorY) {
			return m
		}

		if wall, ok := m.Walls[p]; ok {
			wall.HP--
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			if wall.HP <= 0 {
				delete(m.Walls, p)
			} else {
				m.Walls[p] = wall
			}
			m.Shot = true
			m.ShootMode = false
			m.CurrentScore++
			cur := m.Players[m.CurrentPlayer]
			m.CursorX = cur.X
			m.CursorY = cur.Y
			return m
		}

		if !m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY) {
			var hit bool
			for i, pl := range m.Players {
				if i != m.CurrentPlayer && pl.X == m.CursorX && pl.Y == m.CursorY {
					m.Players[i].HP--
					m.BoxTrigger = TriggerDamage
					m.TriggerTimer = 6
					if m.Players[i].HP <= 0 {
						m.CurrentScore -= 5
						if i < m.CurrentPlayer {
							m.CurrentPlayer--
						}
						m.Players = append(m.Players[:i], m.Players[i+1:]...)
						if len(m.Players) == 0 {
							return m
						}
						if m.CurrentPlayer >= len(m.Players) {
							m.CurrentPlayer = 0
						}
					}
					hit = true
					break
				}
			}
			if !hit {
				for i, en := range m.Enemys {
					if en.X == m.CursorX && en.Y == m.CursorY {
						m.Enemys[i].HP--
						m.BoxTrigger = TriggerDamage
						m.TriggerTimer = 6
						if m.Enemys[i].HP <= 0 {
							m.CurrentScore += 10
							m.Enemys = append(m.Enemys[:i], m.Enemys[i+1:]...)
						}
						hit = true
						break
					}
				}
			}
			if hit {
				m.CurrentScore++
			}
			m.Shot = true
			m.ShootMode = false
			cur := m.Players[m.CurrentPlayer]
			m.CursorX = cur.X
			m.CursorY = cur.Y
		}
	} else if !m.ShootMode && !m.UltMode && !m.PushStrikeMode && !m.Moved {
		if m.IsInRange(m.CursorX, m.CursorY) && !m.IsWall(p) && !wallBlocked && !m.OccupiedByOther(m.CursorX, m.CursorY) {
			m.Players[m.CurrentPlayer].X = m.CursorX
			m.Players[m.CurrentPlayer].Y = m.CursorY

			if m.Water[p] {
				m.Players[m.CurrentPlayer].Effects = ResolveEffects(
					m.Players[m.CurrentPlayer].Effects,
					Effect{Type: EffectWet, Duration: 2},
				)
			}
			if m.FireTiles[p] > 0 {
				if !HasEffect(m.Players[m.CurrentPlayer].Effects, EffectWet) {
					m.Players[m.CurrentPlayer].Effects = ResolveEffects(
						m.Players[m.CurrentPlayer].Effects,
						Effect{Type: EffectFire, Duration: 2},
					)
				}
			}
			if m.SmokeTiles[p] > 0 {
				m.Players[m.CurrentPlayer].Effects = ResolveEffects(
					m.Players[m.CurrentPlayer].Effects,
					Effect{Type: EffectSmoke, Duration: 2},
				)
			}

			m.Moved = true
			m.CursorX = m.Players[m.CurrentPlayer].X
			m.CursorY = m.Players[m.CurrentPlayer].Y
		}
	}

	return m
}

func (m *Model) doUlt() *Model {
	if len(m.Players) == 0 || m.CurrentPlayer >= len(m.Players) {
		return m
	}
	current := m.Players[m.CurrentPlayer]
	if current.UltCharges <= 0 {
		return m
	}
	if !m.ultInAxisRange(m.CursorX, m.CursorY) {
		return m
	}

	m.Players[m.CurrentPlayer].UltCharges--
	m.UltMode = false
	m.UltAxis = ""
	m.Shot = true
	m.CurrentScore += 2

	affected := m.ultCross(m.CursorX, m.CursorY)

	affectedSet := make(map[Point]bool, len(affected))
	for _, p := range affected {
		affectedSet[p] = true
	}

	for _, p := range affected {
		if m.Water[p] || m.SmokeTiles[p] > 0 {
			m.SmokeTiles[p] = 2
		} else {
			m.FireTiles[p] = 2
		}
	}

	for i, pl := range m.Players {
		p := Point{pl.X, pl.Y}
		if !affectedSet[p] {
			continue
		}
		if m.SmokeTiles[p] > 0 {
			m.Players[i].Effects = ResolveEffects(
				m.Players[i].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		} else if m.FireTiles[p] > 0 {
			if HasEffect(pl.Effects, EffectWet) {
				m.Players[i].Effects = RemoveEffect(m.Players[i].Effects, EffectWet)
			} else {
				m.Players[i].Effects = ResolveEffects(
					m.Players[i].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
		}
	}

	for i, en := range m.Enemys {
		p := Point{en.X, en.Y}
		if !affectedSet[p] {
			continue
		}
		if m.SmokeTiles[p] > 0 {
			m.Enemys[i].Effects = ResolveEffects(
				m.Enemys[i].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		} else if m.FireTiles[p] > 0 {
			if HasEffect(en.Effects, EffectWet) {
				m.Enemys[i].Effects = RemoveEffect(m.Enemys[i].Effects, EffectWet)
			} else {
				m.Enemys[i].Effects = ResolveEffects(
					m.Enemys[i].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
		}
	}

	return m
}

func (m *Model) doPushStrike() *Model {
	if len(m.Players) == 0 || m.CurrentPlayer >= len(m.Players) {
		return m
	}
	current := m.Players[m.CurrentPlayer]

	if HasEffect(current.Effects, EffectSmoke) {
		m.PushStrikeMode = false
		m.Shot = true
		return m
	}

	if !m.ultInAxisRange(m.CursorX, m.CursorY) {
		return m
	}

	m.PushStrikeMode = false
	m.UltAxis = ""
	m.Shot = true

	cx, cy := m.CursorX, m.CursorY
	center := Point{cx, cy}

	if wall, ok := m.Walls[center]; ok {
		wall.HP--
		m.BoxTrigger = TriggerDamage
		m.TriggerTimer = 6
		if wall.HP <= 0 {
			delete(m.Walls, center)
		} else {
			m.Walls[center] = wall
		}
		m.CurrentScore++
	} else {
		var hit bool
		for i, pl := range m.Players {
			if pl.X == cx && pl.Y == cy {
				m.Players[i].HP--
				m.BoxTrigger = TriggerDamage
				m.TriggerTimer = 6
				if m.Players[i].HP <= 0 {
					m.CurrentScore -= 5
					if i < m.CurrentPlayer {
						m.CurrentPlayer--
					}
					m.Players = append(m.Players[:i], m.Players[i+1:]...)
					if len(m.Players) == 0 {
						return m
					}
					if m.CurrentPlayer >= len(m.Players) {
						m.CurrentPlayer = 0
					}
				}
				hit = true
				break
			}
		}
		if !hit {
			for i, en := range m.Enemys {
				if en.X == cx && en.Y == cy {
					m.Enemys[i].HP--
					m.BoxTrigger = TriggerDamage
					m.TriggerTimer = 6
					if m.Enemys[i].HP <= 0 {
						m.CurrentScore += 10
						m.Enemys = append(m.Enemys[:i], m.Enemys[i+1:]...)
					}
					hit = true
					break
				}
			}
		}
		if hit {
			m.CurrentScore++
		}
	}

	pushOffsets := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, po := range pushOffsets {
		src := Point{cx + po.X, cy + po.Y}
		dst := Point{cx + 2*po.X, cy + 2*po.Y}

		if dst.X < 0 || dst.X >= GridW || dst.Y < 0 || dst.Y >= GridH {
			continue
		}

		pushed := false
		for i, pl := range m.Players {
			if pl.X == src.X && pl.Y == src.Y {
				m.Players[i].X = dst.X
				m.Players[i].Y = dst.Y
				pushed = true
				break
			}
		}
		if !pushed {
			for i, en := range m.Enemys {
				if en.X == src.X && en.Y == src.Y {
					m.Enemys[i].X = dst.X
					m.Enemys[i].Y = dst.Y
					break
				}
			}
		}
	}

	cur := m.Players[m.CurrentPlayer]
	m.CursorX = cur.X
	m.CursorY = cur.Y
	return m
}
