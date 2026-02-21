package tests

import (
	"testing"

	"hera/generate"
	"hera/utils"

	"github.com/stretchr/testify/assert"
	"github.com/charmbracelet/lipgloss"
)

// --- clamp ---

func TestClamp(t *testing.T) {
	assert.Equal(t, 0, utils.Clamp(-5, 0, 9))
	assert.Equal(t, 9, utils.Clamp(15, 0, 9))
	assert.Equal(t, 5, utils.Clamp(5, 0, 9))
	assert.Equal(t, 0, utils.Clamp(0, 0, 9))
	assert.Equal(t, 9, utils.Clamp(9, 0, 9))
}

// --- abs ---

func TestAbs(t *testing.T) {
	assert.Equal(t, 5, utils.Abs(-5))
	assert.Equal(t, 5, utils.Abs(5))
	assert.Equal(t, 0, utils.Abs(0))
}

// --- generateTiles ---

func TestGenerateTiles_Count(t *testing.T) {
	tiles := generate.GenerateTiles(5, 5, 10, nil)
	assert.Len(t, tiles, 10)
}

func TestGenerateTiles_NotNearPlayer(t *testing.T) {
	playerX, playerY := 5, 5
	tiles := generate.GenerateTiles(playerX, playerY, 20, nil)

	for p := range tiles {
		assert.False(t,
			utils.Abs(p.X-playerX) <= 1 && utils.Abs(p.Y-playerY) <= 1,
			"Tile close to the player: %v ", p,
		)
	}
}

func TestGenerateTiles_NoOverlapWithBlocked(t *testing.T) {
	walls := generate.GenerateTiles(5, 5, 10, nil)
	water := generate.GenerateTiles(5, 5, 10, walls)

	for p := range water {
		assert.False(t, walls[p], "Water and wall overlap: %v", p)
	}
}

// --- Move ---

func testModel() generate.Model {
    walls := map[generate.Point]bool{
        {X: 3, Y: 5}: true,
    }
    water := map[generate.Point]bool{
        {X: 5, Y: 3}: true,
    }
    players := []generate.Player{
        {X: 4, Y: 5, Style: lipgloss.NewStyle()},
        {X: 9, Y: 9, Style: lipgloss.NewStyle()},
    }
    return generate.Model{
        Players:       players,
        CurrentPlayer: 0,
        CursorX:       4,
        CursorY:       5,
        Walls:         walls,
        Water:         water,
    }
}

func TestMove_Normal(t *testing.T) {
    m := testModel()
    m.CursorX, m.CursorY = 4, 4
    p := generate.Point{X: m.CursorX, Y: m.CursorY}
    if !m.Walls[p] && !m.Water[p] {
        m.Players[m.CurrentPlayer].X = m.CursorX
        m.Players[m.CurrentPlayer].Y = m.CursorY
    }
    assert.Equal(t, 4, m.Players[0].X)
    assert.Equal(t, 4, m.Players[0].Y)
}

func TestMove_BlockedByWall(t *testing.T) {
    m := testModel()
    m.CursorX, m.CursorY = 3, 5
    p := generate.Point{X: m.CursorX, Y: m.CursorY}
    if !m.Walls[p] && !m.Water[p] {
        m.Players[m.CurrentPlayer].X = m.CursorX
        m.Players[m.CurrentPlayer].Y = m.CursorY
    }
    assert.Equal(t, 4, m.Players[0].X)
    assert.Equal(t, 5, m.Players[0].Y)
}

func TestMove_BlockedByWater(t *testing.T) {
    m := testModel()
    m.CursorX, m.CursorY = 5, 3
    p := generate.Point{X: m.CursorX, Y: m.CursorY}
    if !m.Walls[p] && !m.Water[p] {
        m.Players[m.CurrentPlayer].X = m.CursorX
        m.Players[m.CurrentPlayer].Y = m.CursorY
    }
    assert.Equal(t, 4, m.Players[0].X)
    assert.Equal(t, 5, m.Players[0].Y)
}

func TestMove_ClampedAtBorder(t *testing.T) {
    m := testModel()
    m.Players[0].X, m.Players[0].Y = 0, 0
    m.CursorX = utils.Clamp(-1, 0, generate.GridW-1)
    m.CursorY = utils.Clamp(0, 0, generate.GridH-1)
    assert.Equal(t, 0, m.CursorX)
    assert.Equal(t, 0, m.CursorY)

    m.Players[0].X, m.Players[0].Y = generate.GridW-1, generate.GridH-1
    m.CursorX = utils.Clamp(generate.GridW, 0, generate.GridW-1)
    assert.Equal(t, generate.GridW-1, m.CursorX)
}

func TestTurnAdvances(t *testing.T) {
    m := testModel()
    assert.Equal(t, 0, m.CurrentPlayer)
    m.CurrentPlayer = (m.CurrentPlayer + 1) % len(m.Players)
    assert.Equal(t, 1, m.CurrentPlayer)
}

func TestOccupiedByOther(t *testing.T) {
    m := testModel()
    assert.True(t, m.OccupiedByOther(9, 9))
    assert.False(t, m.OccupiedByOther(6, 6))
}
