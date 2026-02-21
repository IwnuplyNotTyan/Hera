package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- clamp ---

func TestClamp(t *testing.T) {
	assert.Equal(t, 0, clamp(-5, 0, 9))
	assert.Equal(t, 9, clamp(15, 0, 9))
	assert.Equal(t, 5, clamp(5, 0, 9))
	assert.Equal(t, 0, clamp(0, 0, 9))
	assert.Equal(t, 9, clamp(9, 0, 9))
}

// --- abs ---

func TestAbs(t *testing.T) {
	assert.Equal(t, 5, abs(-5))
	assert.Equal(t, 5, abs(5))
	assert.Equal(t, 0, abs(0))
}

// --- generateTiles ---

func TestGenerateTiles_Count(t *testing.T) {
	tiles := generateTiles(5, 5, 10, nil)
	assert.Len(t, tiles, 10)
}

func TestGenerateTiles_NotNearPlayer(t *testing.T) {
	playerX, playerY := 5, 5
	tiles := generateTiles(playerX, playerY, 20, nil)

	for p := range tiles {
		assert.False(t,
			abs(p.x-playerX) <= 1 && abs(p.y-playerY) <= 1,
			"Tile close to the player: %v ", p,
		)
	}
}

func TestGenerateTiles_NoOverlapWithBlocked(t *testing.T) {
	walls := generateTiles(5, 5, 10, nil)
	water := generateTiles(5, 5, 10, walls)

	for p := range water {
		assert.False(t, walls[p], "Water and wall overlap: %v", p)
	}
}

// --- Move ---

func testModel() model {
	walls := map[point]bool{
		{3, 5}: true,
	}
	water := map[point]bool{
		{5, 3}: true,
	}
	return model{
		x:     4,
		y:     5,
		walls: walls,
		water: water,
	}
}

func TestMove_Normal(t *testing.T) {
	m := testModel()
	m = m.Move(4, 4)
	assert.Equal(t, 4, m.x)
	assert.Equal(t, 4, m.y)
}

func TestMove_BlockedByWall(t *testing.T) {
	m := testModel()
	m = m.Move(3, 5)
	assert.Equal(t, 4, m.x)
	assert.Equal(t, 5, m.y)
}

func TestMove_BlockedByWater(t *testing.T) {
	m := testModel()
	m = m.Move(5, 3)
	assert.Equal(t, 4, m.x)
	assert.Equal(t, 5, m.y)
}

func TestMove_ClampedAtBorder(t *testing.T) {
	m := testModel()
	m.x, m.y = 0, 0
	m = m.Move(-1, 0)
	assert.Equal(t, 0, m.x)
	assert.Equal(t, 0, m.y)

	m.x, m.y = gridW-1, gridH-1
	m = m.Move(gridW, gridH-1)
	assert.Equal(t, gridW-1, m.x)
	assert.Equal(t, gridH-1, m.y)
}
