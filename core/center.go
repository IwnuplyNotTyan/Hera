package generate

func centerToGrid(mode string) (int, int) {
	switch mode {
	case "tl":
		return 0, 0
	case "tc":
		return 0, 1
	case "tr":
		return 0, 2
	case "cl":
		return 1, 0
	case "c":
		return 1, 1
	case "cr":
		return 1, 2
	case "bl":
		return 2, 0
	case "bc":
		return 2, 1
	case "br":
		return 2, 2
	default:
		return 1, 1
	}
}

func gridToCenter(row, col int) string {
	grid := [3][3]string{
		{"tl", "tc", "tr"},
		{"cl", "c", "cr"},
		{"bl", "bc", "br"},
	}
	return grid[row][col]
}
