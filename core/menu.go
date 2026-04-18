package generate

import (
	bz "github.com/lrstanley/bubblezone"
)

func (m *Model) startGame() {
	players := []Player{}
	starts := []Point{
		{X: 1, Y: 1},
		{X: GridW - 2, Y: GridH - 2},
		{X: GridW - 2, Y: 1},
		{X: 1, Y: GridH - 2},
	}

	playerCount := 2
	if playerCount > 4 {
		playerCount = 4
	}

	for i := 0; i < playerCount; i++ {
		players = append(players, Player{
			X:          starts[i].X,
			Y:          starts[i].Y,
			HP:         MaxHP,
			UltCharges: maxUltCharges,
			Style:      m.Styles.PlayerStyles[i%len(m.Styles.EnemysStyles)],
		})
	}

	blocked := make(map[Point]bool)
	for _, p := range players {
		blocked[Point{p.X, p.Y}] = true
	}
	for p := range m.Walls {
		blocked[p] = true
	}
	for p := range m.Water {
		blocked[p] = true
	}

	walls := GenerateTiles(GridW/2, GridH/2, wallCount, nil)
	water := GenerateTiles(GridW/2, GridH/2, waterCount, walls)

	m.Players = players
	m.Enemys = []Enemy{}
	m.Walls = walls
	m.Water = water
	m.FireTiles = make(map[Point]int)
	m.SmokeTiles = make(map[Point]int)
	m.CurrentPlayer = 0
	m.CurrentEnemy = 0
	m.CursorX = players[0].X
	m.CursorY = players[0].Y
	m.Moved = false
	m.Shot = false
	m.ShootMode = false
	m.UltMode = false
	m.UltAxis = ""
	m.EnemyTurn = false
	m.EnemyIdx = 0
	m.MenuSelected = 0
	m.Z = bz.New()
}

func (m *Model) nextTheme() {
	themes := m.AvailableThemes
	if len(themes) == 0 {
		themes = []string{"default"}
	}
	currentIdx := 0
	for i, t := range themes {
		if t == m.ThemeName {
			currentIdx = i
			break
		}
	}
	nextIdx := (currentIdx + 1) % len(themes)
	m.ThemeName = themes[nextIdx]
	if m.Theme != nil {
		m.Theme.SetTintID(m.ThemeName)
	}
	m.Styles = NewStyles(m.Theme)
}

func (m *Model) navigateTheme(direction int) {
	currentIdx := 0
	for i, t := range m.AvailableThemes {
		if t == m.ThemeName {
			currentIdx = i
			break
		}
	}
	nextIdx := currentIdx + direction
	if nextIdx < 0 {
		nextIdx = len(m.AvailableThemes) - 1
	} else if nextIdx >= len(m.AvailableThemes) {
		nextIdx = 0
	}
	m.ThemeName = m.AvailableThemes[nextIdx]
	if m.Theme != nil {
		m.Theme.SetTintID(m.ThemeName)
	}
	m.Styles = NewStyles(m.Theme)
}

