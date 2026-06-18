package generate

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) nextTurn() *Model {
	if len(m.Players) == 0 {
		return m
	}
	m.Moved = false
	m.Shot = false
	m.ShootMode = false
	m.UltMode = false
	m.UltAxis = ""

	if HasEffect(m.Players[m.CurrentPlayer].Effects, EffectFire) && m.Players[m.CurrentPlayer].HP > 1 {
		m.Players[m.CurrentPlayer].HP--
		m.BoxTrigger = TriggerDamage
		m.TriggerTimer = 6
	}

	m.Players[m.CurrentPlayer].Effects = TickEffects(
		m.Players[m.CurrentPlayer].Effects,
	)

	p := Point{m.Players[m.CurrentPlayer].X, m.Players[m.CurrentPlayer].Y}
	if m.Water[p] {
		m.Players[m.CurrentPlayer].Effects = ResolveEffects(
			m.Players[m.CurrentPlayer].Effects,
			Effect{Type: EffectWet, Duration: 2},
		)
	}

	if m.CurrentPlayer == len(m.Players)-1 {
		m = m.tickFireTiles()
		for i := range m.Players {
			if i == m.CurrentPlayer {
				continue
			}
			m.Players[i].Effects = TickEffects(m.Players[i].Effects)
		}
	}

	m.CurrentPlayer = (m.CurrentPlayer + 1) % len(m.Players)
	m.CurrentScore++
	next := m.Players[m.CurrentPlayer]
	m.CursorX = next.X
	m.CursorY = next.Y
	return m
}

func (m *Model) tickFireTiles() *Model {
	for p, t := range m.FireTiles {
		t--
		if t <= 0 {
			delete(m.FireTiles, p)
		} else {
			m.FireTiles[p] = t
		}
	}
	for p, t := range m.SmokeTiles {
		t--
		if t <= 0 {
			delete(m.SmokeTiles, p)
		} else {
			m.SmokeTiles[p] = t
		}
	}
	return m
}

func (m *Model) triggerTickCmd() tea.Cmd {
	if m.TriggerTimer > 0 {
		return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
			return triggerTickMsg{}
		})
	}
	return nil
}
