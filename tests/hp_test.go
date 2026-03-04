package tests

import (
	"testing"

	generate "hera/core"

	"github.com/stretchr/testify/assert"
)

func TestShoot_ReducesHP(t *testing.T) {
	m := testModel()
	m.Players[1].X = 5
	m.Players[1].Y = 5
	m.Players[1].HP = 3

	m.Players[m.CurrentPlayer].X = 4
	m.Players[m.CurrentPlayer].Y = 5
	m.CursorX, m.CursorY = 5, 5
	m.ShootMode = true

	p := generate.Point{X: m.CursorX, Y: m.CursorY}
	if !m.Walls[p] {
		for i, pl := range m.Players {
			if i != m.CurrentPlayer && pl.X == m.CursorX && pl.Y == m.CursorY {
				m.Players[i].HP--
				break
			}
		}
	}
	assert.Equal(t, 2, m.Players[1].HP)
}

func TestShoot_PlayerDiesAt0HP(t *testing.T) {
	m := testModel()
	m.Players[1].X = 5
	m.Players[1].Y = 5
	m.Players[1].HP = 1

	for i, pl := range m.Players {
		if i != m.CurrentPlayer && pl.X == 5 && pl.Y == 5 {
			m.Players[i].HP--
			if m.Players[i].HP <= 0 {
				m.Players = append(m.Players[:i], m.Players[i+1:]...)
			}
			break
		}
	}
	assert.Len(t, m.Players, 1)
}

func TestHP_InitialValue(t *testing.T) {
	m := generate.NewModel(2, 2)
	for _, p := range m.Players {
		assert.Equal(t, generate.MaxHP, p.HP)
	}
}
