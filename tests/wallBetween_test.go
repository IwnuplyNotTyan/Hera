package tests

import (
	generate "hera/core"
	"testing"

	"github.com/charmbracelet/lipgloss"
	bubbletint "github.com/lrstanley/bubbletint"
	"github.com/stretchr/testify/assert"
)

func TestHasWallBetween_Clear(t *testing.T) {
	m := testModel()
	assert.False(t, m.HasWallBetweenPoints(4, 5, 4, 4))
}

func TestHasWallBetween_WallBlocks(t *testing.T) {
	m := testModel()
	assert.True(t, m.HasWallBetweenPoints(4, 5, 2, 5))
}

func TestHasWallBetween_StartNotCounted(t *testing.T) {
	bubbletint.NewDefaultRegistry()
	generate.RegisterDefaults()
	theme := bubbletint.DefaultRegistry
	walls := map[generate.Point]generate.Wall{
		{X: 4, Y: 5}: {HP: generate.WallHP},
	}
	players := []generate.Player{
		{X: 4, Y: 5, HP: generate.MaxHP, Style: lipgloss.NewStyle()},
		{X: 9, Y: 9, HP: generate.MaxHP, Style: lipgloss.NewStyle()},
	}
	m := generate.Model{
		Theme:         theme,
		Styles:        generate.NewStyles(theme),
		Players:       players,
		CurrentPlayer: 0,
		CursorX:       4, CursorY: 5,
		Walls:      walls,
		Water:      map[generate.Point]bool{},
		FireTiles:  map[generate.Point]int{},
		SmokeTiles: map[generate.Point]int{},
		Enemys:     []generate.Enemy{},
	}
	assert.False(t, m.HasWallBetweenPoints(4, 5, 4, 4))
}
