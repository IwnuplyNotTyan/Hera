package generate

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"hera/utils"
)

func (m *Model) closestPlayer(ex, ey int) (int, int) {
	if len(m.Players) == 0 {
		return ex, ey
	}
	bestX, bestY := m.Players[0].X, m.Players[0].Y
	bestDist := utils.Abs(ex-bestX) + utils.Abs(ey-bestY)
	for _, pl := range m.Players[1:] {
		d := utils.Abs(ex-pl.X) + utils.Abs(ey-pl.Y)
		if d < bestDist {
			bestDist = d
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
		tx, ty := m.closestPlayer(en.X, en.Y)
		dist := utils.Abs(en.X-tx) + utils.Abs(en.Y-ty)

		if dist <= shootRange && !m.HasWallBetweenPoints(en.X, en.Y, tx, ty) {
			for j, pl := range m.Players {
				if pl.X == tx && pl.Y == ty {
					m.Players[j].HP--
					m.BoxTrigger = TriggerDamage
					m.TriggerTimer = 6
					if m.Players[j].HP <= 0 {
						if j < m.CurrentPlayer {
							m.CurrentPlayer--
						}
						m.Players = append(m.Players[:j], m.Players[j+1:]...)
						if len(m.Players) == 0 {
							return m
						}
						if m.CurrentPlayer >= len(m.Players) {
							m.CurrentPlayer = 0
						}
					}
					break
				}
			}
			return m
		}

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

		moved := false
		for _, mv := range moves {
			if mv.X < 0 || mv.X >= GridW || mv.Y < 0 || mv.Y >= GridH {
				continue
			}
			if m.Walls[mv] || m.enemyOccupied(mv.X, mv.Y, idx) {
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

			moved = true
			break
		}
		if !moved {
			break
		}
	}

	if HasEffect(m.Enemys[idx].Effects, EffectFire) && m.Enemys[idx].HP > 1 {
		m.Enemys[idx].HP--
	}

	m.Enemys[idx].Effects = TickEffects(m.Enemys[idx].Effects)
	return m
}

func enemyTurnCmd(idx int) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return enemyTurnMsg{enemyIdx: idx}
	})
}
