package generate

import "hera/utils"

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
	} else if m.RamMode && !m.Shot {
		m = m.doRam()
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
	} else if !m.ShootMode && !m.UltMode && !m.PushStrikeMode && !m.RamMode && !m.Moved {
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

		var srcPl, srcEn = -1, -1
		for i, pl := range m.Players {
			if pl.X == src.X && pl.Y == src.Y {
				srcPl = i
				break
			}
		}
		if srcPl == -1 {
			for i, en := range m.Enemys {
				if en.X == src.X && en.Y == src.Y {
					srcEn = i
					break
				}
			}
		}
		if srcPl == -1 && srcEn == -1 {
			continue
		}

		var dstPl, dstEn = -1, -1
		for i, pl := range m.Players {
			if pl.X == dst.X && pl.Y == dst.Y {
				dstPl = i
				break
			}
		}
		if dstPl == -1 {
			for i, en := range m.Enemys {
				if en.X == dst.X && en.Y == dst.Y {
					dstEn = i
					break
				}
			}
		}

		if dstPl >= 0 {
			m.Players[dstPl].HP--
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			if m.Players[dstPl].HP <= 0 {
				m.CurrentScore -= 5
				if dstPl < m.CurrentPlayer {
					m.CurrentPlayer--
				}
				m.Players = append(m.Players[:dstPl], m.Players[dstPl+1:]...)
				if len(m.Players) == 0 {
					return m
				}
				if m.CurrentPlayer >= len(m.Players) {
					m.CurrentPlayer = 0
				}
				if srcPl > dstPl {
					srcPl--
				}
			}
			m.CurrentScore++
		} else if dstEn >= 0 {
			m.Enemys[dstEn].HP--
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			if m.Enemys[dstEn].HP <= 0 {
				m.CurrentScore += 10
				m.Enemys = append(m.Enemys[:dstEn], m.Enemys[dstEn+1:]...)
				if srcEn > dstEn {
					srcEn--
				}
			} else {
				m.CurrentScore++
			}
		}

		if dstPl >= 0 || dstEn >= 0 {
			if srcPl >= 0 {
				m.Players[srcPl].HP--
				m.BoxTrigger = TriggerDamage
				m.TriggerTimer = 6
				if m.Players[srcPl].HP <= 0 {
					m.CurrentScore -= 5
					if srcPl < m.CurrentPlayer {
						m.CurrentPlayer--
					}
					m.Players = append(m.Players[:srcPl], m.Players[srcPl+1:]...)
					if len(m.Players) == 0 {
						return m
					}
					if m.CurrentPlayer >= len(m.Players) {
						m.CurrentPlayer = 0
					}
				} else {
					m.CurrentScore++
				}
			} else {
				m.Enemys[srcEn].HP--
				m.BoxTrigger = TriggerDamage
				m.TriggerTimer = 6
				if m.Enemys[srcEn].HP <= 0 {
					m.CurrentScore += 10
					m.Enemys = append(m.Enemys[:srcEn], m.Enemys[srcEn+1:]...)
				} else {
					m.CurrentScore++
				}
			}
			continue
		}

		if srcPl >= 0 {
			m.Players[srcPl].X = dst.X
			m.Players[srcPl].Y = dst.Y
		} else {
			m.Enemys[srcEn].X = dst.X
			m.Enemys[srcEn].Y = dst.Y
		}
	}

	cur := m.Players[m.CurrentPlayer]
	m.CursorX = cur.X
	m.CursorY = cur.Y
	return m
}

func (m *Model) doRam() *Model {
	if len(m.Players) == 0 || m.CurrentPlayer >= len(m.Players) {
		return m
	}
	cur := m.Players[m.CurrentPlayer]

	if HasEffect(cur.Effects, EffectSmoke) {
		m.RamMode = false
		m.Shot = true
		return m
	}

	if !m.ultInAxisRange(m.CursorX, m.CursorY) {
		return m
	}

	cx, cy := cur.X, cur.Y
	tx, ty := m.CursorX, m.CursorY

	dx := 0
	if tx > cx {
		dx = 1
	} else if tx < cx {
		dx = -1
	}
	dy := 0
	if ty > cy {
		dy = 1
	} else if ty < cy {
		dy = -1
	}

	steps := utils.Abs(tx-cx) + utils.Abs(ty-cy)

	m.RamMode = false
	m.UltAxis = ""
	m.Shot = true
	m.CurrentScore += 2

	for step := 1; step <= steps; step++ {
		cell := Point{cx + dx*step, cy + dy*step}

		if wall, ok := m.Walls[cell]; ok {
			wall.HP--
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			m.CurrentScore++
			if wall.HP <= 0 {
				delete(m.Walls, cell)
				m.Players[m.CurrentPlayer].X = cell.X
				m.Players[m.CurrentPlayer].Y = cell.Y
			} else {
				m.Walls[cell] = wall
			}
			break
		}

		var plAtCell, enAtCell = -1, -1
		for i := range m.Players {
			if i != m.CurrentPlayer && m.Players[i].X == cell.X && m.Players[i].Y == cell.Y {
				plAtCell = i
				break
			}
		}
		if plAtCell == -1 {
			for i := range m.Enemys {
				if m.Enemys[i].X == cell.X && m.Enemys[i].Y == cell.Y {
					enAtCell = i
					break
				}
			}
		}

		if plAtCell >= 0 || enAtCell >= 0 {
			pushDst := Point{cell.X + dx, cell.Y + dy}
			blocked := pushDst.X < 0 || pushDst.X >= GridW || pushDst.Y < 0 || pushDst.Y >= GridH
			if !blocked {
				if m.IsWall(pushDst) {
					blocked = true
				} else if m.OccupiedByOther(pushDst.X, pushDst.Y) {
					blocked = true
				}
			}

			if !blocked {
				if plAtCell >= 0 {
					m.Players[plAtCell].X = pushDst.X
					m.Players[plAtCell].Y = pushDst.Y
				} else {
					m.Enemys[enAtCell].X = pushDst.X
					m.Enemys[enAtCell].Y = pushDst.Y
				}
				m.Players[m.CurrentPlayer].X = cell.X
				m.Players[m.CurrentPlayer].Y = cell.Y
			}

			dmg := 2
			if blocked {
				dmg = 3
			}

			if plAtCell >= 0 {
				m.Players[plAtCell].HP -= dmg
				m.BoxTrigger = TriggerDamage
				m.TriggerTimer = 6
				if m.Players[plAtCell].HP <= 0 {
					m.CurrentScore -= 5
					if plAtCell < m.CurrentPlayer {
						m.CurrentPlayer--
					}
					m.Players = append(m.Players[:plAtCell], m.Players[plAtCell+1:]...)
					if len(m.Players) == 0 {
						return m
					}
					if m.CurrentPlayer >= len(m.Players) {
						m.CurrentPlayer = 0
					}
				}
			} else {
				m.Enemys[enAtCell].HP -= dmg
				m.BoxTrigger = TriggerDamage
				m.TriggerTimer = 6
				if m.Enemys[enAtCell].HP <= 0 {
					m.CurrentScore += 10
					m.Enemys = append(m.Enemys[:enAtCell], m.Enemys[enAtCell+1:]...)
				}
			}
			break
		}

		m.Players[m.CurrentPlayer].X = cell.X
		m.Players[m.CurrentPlayer].Y = cell.Y
		cur = m.Players[m.CurrentPlayer]

		if m.Water[cell] {
			m.Players[m.CurrentPlayer].Effects = ResolveEffects(
				m.Players[m.CurrentPlayer].Effects,
				Effect{Type: EffectWet, Duration: 2},
			)
		}
		if m.FireTiles[cell] > 0 {
			if !HasEffect(m.Players[m.CurrentPlayer].Effects, EffectWet) {
				m.Players[m.CurrentPlayer].Effects = ResolveEffects(
					m.Players[m.CurrentPlayer].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
		}
		if m.SmokeTiles[cell] > 0 {
			m.Players[m.CurrentPlayer].Effects = ResolveEffects(
				m.Players[m.CurrentPlayer].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		}
	}

	if m.CurrentPlayer < len(m.Players) {
		m.Players[m.CurrentPlayer].HP--
		m.BoxTrigger = TriggerDamage
		m.TriggerTimer = 6
		if m.Players[m.CurrentPlayer].HP <= 0 {
			m.CurrentScore -= 5
			m.Players = append(m.Players[:m.CurrentPlayer], m.Players[m.CurrentPlayer+1:]...)
			if len(m.Players) == 0 {
				return m
			}
			if m.CurrentPlayer >= len(m.Players) {
				m.CurrentPlayer = 0
			}
		}
	}

	if m.CurrentPlayer < len(m.Players) {
		cur = m.Players[m.CurrentPlayer]
		m.CursorX = cur.X
		m.CursorY = cur.Y
	}
	return m
}
