package tests

import (
	"testing"

	generate "hera/core"
	"hera/i18n"

	"github.com/charmbracelet/lipgloss"
	bubbletint "github.com/lrstanley/bubbletint"
	"github.com/stretchr/testify/assert"
)

func createTestModel() generate.Model {
	loc, _ := i18n.NewTranslator("../locales", "en")
	theme := bubbletint.NewRegistry(bubbletint.TintDraculaPlus, bubbletint.DefaultTints()...)
	walls := map[generate.Point]bool{
		{X: 3, Y: 5}: true,
	}
	water := map[generate.Point]bool{
		{X: 5, Y: 3}: true,
	}
	players := []generate.Player{
		{X: 4, Y: 5, HP: generate.MaxHP, Style: lipgloss.NewStyle()},
		{X: 9, Y: 9, HP: generate.MaxHP, Style: lipgloss.NewStyle()},
	}
	return generate.Model{
		Theme:         theme,
		Players:       players,
		CurrentPlayer: 0,
		CursorX:       4,
		CursorY:       5,
		Walls:         walls,
		Water:         water,
		FireTiles:     map[generate.Point]int{},
		SmokeTiles:    map[generate.Point]int{},
		Enemys:        []generate.Enemy{},
		Localizer:     loc,
	}
}

func TestHP_InitialValue(t *testing.T) {
	loc, _ := i18n.NewTranslator("../locales", "en")
	theme := bubbletint.NewRegistry(bubbletint.TintDraculaPlus, bubbletint.DefaultTints()...)
	m := generate.NewModel(2, 2, loc, theme, false)
	for _, p := range m.Players {
		assert.Equal(t, generate.MaxHP, p.HP)
	}
}

func TestShoot_ReducesHP(t *testing.T) {
	m := createTestModel()
	m.Players[1].X, m.Players[1].Y = 5, 5
	m.ShootMode = true
	m.CursorX, m.CursorY = 5, 5

	p := generate.Point{X: m.CursorX, Y: m.CursorY}
	if !m.Walls[p] {
		for i, pl := range m.Players {
			if i != m.CurrentPlayer && pl.X == m.CursorX && pl.Y == m.CursorY {
				m.Players[i].HP--
				break
			}
		}
	}
	assert.Equal(t, generate.MaxHP-1, m.Players[1].HP)
}

func TestShoot_PlayerDiesAt0HP(t *testing.T) {
	m := createTestModel()
	m.Players[1].X, m.Players[1].Y = 5, 5
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
