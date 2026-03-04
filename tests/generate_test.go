package tests

import (
	"testing"

	generate "hera/core"
	"hera/utils"

	"github.com/stretchr/testify/assert"
)

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
