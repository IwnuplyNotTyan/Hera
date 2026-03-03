package generate

import (
	"math/rand"

	"hera/utils"
)

func GenerateTiles(playerX, playerY, count int, blocked map[Point]bool) map[Point]bool {
	tiles := make(map[Point]bool)
	for len(tiles) < count {
		x := rand.Intn(GridW)
		y := rand.Intn(GridH)
		p := Point{x, y}

		if utils.Abs(x-playerX) <= 1 && utils.Abs(y-playerY) <= 1 {
			continue
		}
		if blocked[p] {
			continue
		}
		tiles[p] = true
	}
	return tiles
}
