package generate

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"hera/utils"
)

func (m *Model) targetEnemyPlayer(ex, ey int, atk AttackType) (int, int) {
	if len(m.Players) == 0 {
		return ex, ey
	}
	bestX, bestY := m.Players[0].X, m.Players[0].Y
	bestScore := -999
	for _, pl := range m.Players {
		score := -(utils.Abs(ex-pl.X) + utils.Abs(ey-pl.Y)) * 2
		switch atk {
		case AttackRam:
			if ex == pl.X || ey == pl.Y {
				score += 5
			}
		case AttackMeleePush:
			if utils.Abs(ex-pl.X)+utils.Abs(ey-pl.Y) == 1 {
				score += 10
			}
		case AttackPushStrike:
			if ex == pl.X || ey == pl.Y {
				score += 3
			}
		}
		if score > bestScore {
			bestScore = score
			bestX, bestY = pl.X, pl.Y
		}
	}
	return bestX, bestY
}

func (m *Model) enemyOccupied(x, y, skipIdx int) bool {
	for i, e := range m.Enemys {
		if i != skipIdx && e.X == x && e.Y == y {
			return true
		}
	}
	for _, p := range m.Players {
		if p.X == x && p.Y == y {
			return true
		}
	}
	return false
}

func (m *Model) doEnemyTurn(idx int) *Model {
	if len(m.Players) == 0 || idx >= len(m.Enemys) {
		return m
	}

	for step := 0; step < moveRange; step++ {
		en := m.Enemys[idx]
		if m.tryEnemyAttack(idx) {
			goto done
		}

		tx, ty := m.targetEnemyPlayer(en.X, en.Y, en.AttackType)
		if !m.moveEnemyToward(idx, tx, ty) {
			break
		}
	}

	m.tryEnemyAttack(idx)

done:
	if idx < len(m.Enemys) {
		if HasEffect(m.Enemys[idx].Effects, EffectFire) && m.Enemys[idx].HP > 1 {
			m.Enemys[idx].HP--
		}
		m.Enemys[idx].Effects = TickEffects(m.Enemys[idx].Effects)
	}
	return m
}

func (m *Model) moveEnemyToward(idx, tx, ty int) bool {
	en := m.Enemys[idx]
	moves := []Point{}
	if tx > en.X {
		moves = append(moves, Point{en.X + 1, en.Y})
	}
	if tx < en.X {
		moves = append(moves, Point{en.X - 1, en.Y})
	}
	if ty > en.Y {
		moves = append(moves, Point{en.X, en.Y + 1})
	}
	if ty < en.Y {
		moves = append(moves, Point{en.X, en.Y - 1})
	}

	for _, mv := range moves {
		if mv.X < 0 || mv.X >= GridW || mv.Y < 0 || mv.Y >= GridH {
			continue
		}
		if m.IsWall(mv) || m.enemyOccupied(mv.X, mv.Y, idx) {
			continue
		}
		m.Enemys[idx].X = mv.X
		m.Enemys[idx].Y = mv.Y

		p := Point{mv.X, mv.Y}
		if m.FireTiles[p] > 0 && !HasEffect(m.Enemys[idx].Effects, EffectWet) {
			m.Enemys[idx].Effects = ResolveEffects(
				m.Enemys[idx].Effects,
				Effect{Type: EffectFire, Duration: 2},
			)
		}
		if m.Water[p] {
			m.Enemys[idx].Effects = ResolveEffects(
				m.Enemys[idx].Effects,
				Effect{Type: EffectWet, Duration: 2},
			)
		}
		if m.SmokeTiles[p] > 0 {
			m.Enemys[idx].Effects = ResolveEffects(
				m.Enemys[idx].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		}
		return true
	}
	return false
}

func (m *Model) tryEnemyAttack(idx int) bool {
	if idx >= len(m.Enemys) {
		return false
	}
	en := m.Enemys[idx]
	if HasEffect(en.Effects, EffectSmoke) {
		return false
	}
	if len(m.Players) == 0 {
		return false
	}

	tx, ty := m.targetEnemyPlayer(en.X, en.Y, en.AttackType)

	switch en.AttackType {
	case AttackShoot:
		return m.enemyShoot(idx, tx, ty)
	case AttackPushStrike:
		return m.enemyPushStrike(idx, tx, ty)
	case AttackRam:
		return m.enemyRam(idx, tx, ty)
	case AttackMeleePush:
		return m.enemyMeleePush(idx, tx, ty)
	}
	return false
}

func (m *Model) enemyShoot(idx, tx, ty int) bool {
	en := m.Enemys[idx]
	dist := utils.Abs(en.X-tx) + utils.Abs(en.Y-ty)
	if dist > shootRange {
		return false
	}
	if m.HasWallBetweenPoints(en.X, en.Y, tx, ty) {
		return false
	}

	for j := range m.Players {
		if m.Players[j].X == tx && m.Players[j].Y == ty {
			m.Players[j].HP--
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			if m.Players[j].HP <= 0 {
				if j < m.CurrentPlayer {
					m.CurrentPlayer--
				}
				m.Players = append(m.Players[:j], m.Players[j+1:]...)
				if len(m.Players) == 0 {
					return true
				}
				if m.CurrentPlayer >= len(m.Players) {
					m.CurrentPlayer = 0
				}
			}
			return true
		}
	}
	return false
}

func (m *Model) enemyPushStrike(idx, tx, ty int) bool {
	en := m.Enemys[idx]
	if en.X != tx && en.Y != ty {
		return false
	}
	dist := utils.Abs(en.X-tx) + utils.Abs(en.Y-ty)
	if dist < 1 || dist > 4 {
		return false
	}
	if m.HasWallBetweenPoints(en.X, en.Y, tx, ty) {
		return false
	}

	// Center cell at player position — deal 1 damage
	cx, cy := tx, ty

	for j := range m.Players {
		if m.Players[j].X == cx && m.Players[j].Y == cy {
			m.Players[j].HP--
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			if m.Players[j].HP <= 0 {
				if j < m.CurrentPlayer {
					m.CurrentPlayer--
				}
				m.Players = append(m.Players[:j], m.Players[j+1:]...)
				if len(m.Players) == 0 {
					return true
				}
				if m.CurrentPlayer >= len(m.Players) {
					m.CurrentPlayer = 0
				}
			}
			break
		}
	}

	// Push cardinals outward
	pushOffsets := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for _, po := range pushOffsets {
		src := Point{cx + po.X, cy + po.Y}
		dst := Point{cx + 2*po.X, cy + 2*po.Y}

		if dst.X < 0 || dst.X >= GridW || dst.Y < 0 || dst.Y >= GridH {
			continue
		}

		var srcPl, srcEn = -1, -1
		for i := range m.Players {
			if m.Players[i].X == src.X && m.Players[i].Y == src.Y {
				srcPl = i
				break
			}
		}
		if srcPl == -1 {
			for i := range m.Enemys {
				if m.Enemys[i].X == src.X && m.Enemys[i].Y == src.Y {
					srcEn = i
					break
				}
			}
		}
		if srcPl == -1 && srcEn == -1 {
			continue
		}

		blocked := false
		for i := range m.Players {
			if m.Players[i].X == dst.X && m.Players[i].Y == dst.Y {
				blocked = true
				m.Players[i].HP--
				m.BoxTrigger = TriggerDamage
				m.TriggerTimer = 6
				if m.Players[i].HP <= 0 {
					if i < m.CurrentPlayer {
						m.CurrentPlayer--
					}
					m.Players = append(m.Players[:i], m.Players[i+1:]...)
					if len(m.Players) == 0 {
						return true
					}
					if m.CurrentPlayer >= len(m.Players) {
						m.CurrentPlayer = 0
					}
				}
				break
			}
		}
		if !blocked {
			for i := range m.Enemys {
				if m.Enemys[i].X == dst.X && m.Enemys[i].Y == dst.Y {
					blocked = true
					m.Enemys[i].HP--
					if m.Enemys[i].HP <= 0 {
						m.Enemys = append(m.Enemys[:i], m.Enemys[i+1:]...)
					}
					break
				}
			}
		}
		if !blocked {
			if w, ok := m.Walls[dst]; ok {
				w.HP--
				if w.HP <= 0 {
					delete(m.Walls, dst)
				} else {
					m.Walls[dst] = w
				}
				blocked = true
			}
		}

		if blocked {
			if srcPl >= 0 {
				m.Players[srcPl].HP--
				m.BoxTrigger = TriggerDamage
				m.TriggerTimer = 6
				if m.Players[srcPl].HP <= 0 {
					if srcPl < m.CurrentPlayer {
						m.CurrentPlayer--
					}
					m.Players = append(m.Players[:srcPl], m.Players[srcPl+1:]...)
					if len(m.Players) == 0 {
						return true
					}
					if m.CurrentPlayer >= len(m.Players) {
						m.CurrentPlayer = 0
					}
				}
			} else {
				m.Enemys[srcEn].HP--
				if m.Enemys[srcEn].HP <= 0 {
					m.Enemys = append(m.Enemys[:srcEn], m.Enemys[srcEn+1:]...)
				}
			}
		} else {
			if srcPl >= 0 {
				m.Players[srcPl].X = dst.X
				m.Players[srcPl].Y = dst.Y
			} else {
				m.Enemys[srcEn].X = dst.X
				m.Enemys[srcEn].Y = dst.Y
			}
		}
	}

	return true
}

func (m *Model) enemyRam(idx, tx, ty int) bool {
	en := m.Enemys[idx]
	if en.X != tx && en.Y != ty {
		return false
	}
	dist := utils.Abs(en.X-tx) + utils.Abs(en.Y-ty)
	if dist < 2 || dist > 5 {
		return false
	}

	dx, dy := 0, 0
	if tx > en.X {
		dx = 1
	} else if tx < en.X {
		dx = -1
	}
	if ty > en.Y {
		dy = 1
	} else if ty < en.Y {
		dy = -1
	}

	cx, cy := en.X, en.Y

	for step := 1; step <= dist; step++ {
		cell := Point{cx + dx*step, cy + dy*step}

		if wall, ok := m.Walls[cell]; ok {
			wall.HP -= 2
			if wall.HP <= 0 {
				delete(m.Walls, cell)
			} else {
				m.Walls[cell] = wall
			}
			break
		}

		plHit := -1
		for i := range m.Players {
			if m.Players[i].X == cell.X && m.Players[i].Y == cell.Y {
				plHit = i
				break
			}
		}

		if plHit >= 0 {
			m.Players[plHit].HP -= 2
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			if m.Players[plHit].HP <= 0 {
				if plHit < m.CurrentPlayer {
					m.CurrentPlayer--
				}
				m.Players = append(m.Players[:plHit], m.Players[plHit+1:]...)
				if len(m.Players) == 0 {
					return true
				}
				if m.CurrentPlayer >= len(m.Players) {
					m.CurrentPlayer = 0
				}
			}
			break
		}

		pushDst := Point{cell.X + dx, cell.Y + dy}
		enHit := -1
		for i := range m.Enemys {
			if m.Enemys[i].X == cell.X && m.Enemys[i].Y == cell.Y {
				enHit = i
				break
			}
		}

		if enHit >= 0 {
			blocked := pushDst.X < 0 || pushDst.X >= GridW || pushDst.Y < 0 || pushDst.Y >= GridH ||
				m.IsWall(pushDst) || m.enemyOccupied(pushDst.X, pushDst.Y, idx)
			dmg := 3
			if !blocked {
				dmg = 2
				m.Enemys[enHit].X = pushDst.X
				m.Enemys[enHit].Y = pushDst.Y
			}
			m.Enemys[enHit].HP -= dmg
			if m.Enemys[enHit].HP <= 0 {
				m.Enemys = append(m.Enemys[:enHit], m.Enemys[enHit+1:]...)
			}
			break
		}
	}

	return true
}

func (m *Model) enemyMeleePush(idx, tx, ty int) bool {
	en := m.Enemys[idx]
	dist := utils.Abs(en.X-tx) + utils.Abs(en.Y-ty)
	if dist != 1 {
		return false
	}

	pushDst := Point{tx + (tx - en.X), ty + (ty - en.Y)}
	canPush := pushDst.X >= 0 && pushDst.X < GridW && pushDst.Y >= 0 && pushDst.Y < GridH &&
		!m.IsWall(pushDst) && !m.enemyOccupied(pushDst.X, pushDst.Y, idx)

	dmg := 2
	if canPush {
		dmg = 1
	}

	for j := range m.Players {
		if m.Players[j].X == tx && m.Players[j].Y == ty {
			m.Players[j].HP -= dmg
			m.BoxTrigger = TriggerDamage
			m.TriggerTimer = 6
			if m.Players[j].HP <= 0 {
				if j < m.CurrentPlayer {
					m.CurrentPlayer--
				}
				m.Players = append(m.Players[:j], m.Players[j+1:]...)
				if len(m.Players) == 0 {
					return true
				}
				if m.CurrentPlayer >= len(m.Players) {
					m.CurrentPlayer = 0
				}
			} else if canPush {
				m.Players[j].X = pushDst.X
				m.Players[j].Y = pushDst.Y
			}
			return true
		}
	}
	return false
}

func enemyTurnCmd(idx int) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return enemyTurnMsg{enemyIdx: idx}
	})
}
