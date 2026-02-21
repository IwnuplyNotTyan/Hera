package tests

import (
	"testing"

	"hera/generate"
	"hera/utils"

	"github.com/stretchr/testify/assert"
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
	return generate.Model{
		X:     4,
		Y:     5,
		Walls: walls,
		Water: water,
	}
}

func TestMove_Normal(t *testing.T) {
	m := testModel()
	m = m.Move(4, 4)
	assert.Equal(t, 4, m.X)
	assert.Equal(t, 4, m.Y)
}

func TestMove_BlockedByWall(t *testing.T) {
	m := testModel()
	m = m.Move(3, 5)
	assert.Equal(t, 4, m.X)
	assert.Equal(t, 5, m.Y)
}

func TestMove_BlockedByWater(t *testing.T) {
	m := testModel()
	m = m.Move(5, 3)
	assert.Equal(t, 4, m.X)
	assert.Equal(t, 5, m.Y)
}

func TestMove_ClampedAtBorder(t *testing.T) {
	m := testModel()
	m.X, m.Y = 0, 0
	m = m.Move(-1, 0)
	assert.Equal(t, 0, m.X)
	assert.Equal(t, 0, m.Y)

	m.X, m.Y = generate.GridW-1, generate.GridH-1
	m = m.Move(generate.GridW, generate.GridH-1)
	assert.Equal(t, generate.GridW-1, m.X)
	assert.Equal(t, generate.GridH-1, m.Y)
}
