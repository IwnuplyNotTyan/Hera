package tests

import (
	generate "hera/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInRange_FreeCell(t *testing.T) {
	m := testModel()
	assert.True(t, m.IsInRange(4, 4))
	assert.True(t, m.IsInRange(4, 2))
	assert.False(t, m.IsInRange(4, 0))
}

func TestIsInRange_PlayerCellExcluded(t *testing.T) {
	m := testModel()
	assert.False(t, m.IsInRange(4, 5))
}

func TestIsInRange_BlockedByWall(t *testing.T) {
	m := testModel()
	assert.False(t, m.IsInRange(2, 5))
}

func TestReachable_BlockedByWall(t *testing.T) {
	walls := map[generate.Point]bool{
		{X: 3, Y: 3}: true,
		{X: 3, Y: 4}: true,
		{X: 3, Y: 5}: true,
		{X: 3, Y: 6}: true,
		{X: 3, Y: 7}: true,
	}
	players := []generate.Player{
		{X: 4, Y: 5, HP: generate.MaxHP},
		{X: 9, Y: 9, HP: generate.MaxHP},
	}
	m := generate.Model{
		Players:       players,
		CurrentPlayer: 0,
		CursorX:       4, CursorY: 5,
		Walls:      walls,
		Water:      map[generate.Point]bool{},
		FireTiles:  map[generate.Point]int{},
		SteamTiles: map[generate.Point]int{},
		Enemys:     []generate.Enemy{},
	}
	zone := m.Reachable(4, 5, 4)
	assert.False(t, zone[generate.Point{X: 0, Y: 5}], "(0,5) - Unreachable")
	assert.True(t, zone[generate.Point{X: 4, Y: 4}], "(4,4) — Reachable")
	assert.True(t, zone[generate.Point{X: 6, Y: 5}], "(6,5) — Reachable")
}

func TestReachable_ShootRange(t *testing.T) {
	m := testModel()
	m.ShootMode = true
	zone := m.Reachable(4, 5, 2)
	assert.True(t, zone[generate.Point{X: 4, Y: 3}])
	assert.False(t, zone[generate.Point{X: 4, Y: 2}])
}
