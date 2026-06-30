package generate

import (
	"math/rand"
)

func (m *Model) startGame() {
	players := []Player{}
	starts := []Point{
		{X: 1, Y: 1},
		{X: GridW - 2, Y: GridH - 2},
		{X: GridW - 2, Y: 1},
		{X: 1, Y: GridH - 2},
	}

	if m.Seed == 0 {
		m.Seed = RandomSeed()
	}

	rng := rand.New(rand.NewSource(m.Seed))

	intn := func(n int) int {
		return rng.Intn(n)
	}

	playerCount := m.startPlayers
	if playerCount <= 0 {
		playerCount = intn(3) + 2
	}
	enemyCount := m.startEnemies
	if enemyCount <= 0 {
		enemyCount = intn(3) + 2
	}
	if enemyCount > len(m.Styles.EnemysStyles) {
		enemyCount = len(m.Styles.EnemysStyles)
	}

	for i := 0; i < playerCount; i++ {
		effs := make([]Effect, len(m.startPlayerEffects))
		copy(effs, m.startPlayerEffects)
		players = append(players, Player{
			X:          starts[i].X,
			Y:          starts[i].Y,
			HP:         MaxHP,
			UltCharges: maxUltCharges,
			Effects:    effs,
			Style:      m.Styles.PlayerStyles[i%len(m.Styles.PlayerStyles)],
		})
	}

	blocked := make(map[Point]bool)
	for _, p := range players {
		blocked[Point{p.X, p.Y}] = true
	}

	wallPositions := GenerateTiles(GridW/2, GridH/2, wallCount, blocked, rng)
	walls := make(map[Point]Wall, len(wallPositions))
	for p := range wallPositions {
		walls[p] = Wall{HP: WallHP}
		blocked[p] = true
	}

	water := GenerateTiles(GridW/2, GridH/2, waterCount, blocked, rng)
	for p := range water {
		blocked[p] = true
	}

	enemyStarts := GenerateTiles(GridW/2, GridH/2, enemyCount, blocked, rng)
	enemyPositions := make([]Point, 0, enemyCount)
	for p := range enemyStarts {
		enemyPositions = append(enemyPositions, p)
	}

	enemys := make([]Enemy, enemyCount)
	for i := range enemys {
		effs := make([]Effect, len(m.startEnemyEffects))
		copy(effs, m.startEnemyEffects)
		enemys[i] = Enemy{
			X:       enemyPositions[i].X,
			Y:       enemyPositions[i].Y,
			HP:      MaxHP,
			Effects: effs,
			Style:   m.Styles.EnemysStyles[i],
		}
	}

	m.Players = players
	m.Enemys = enemys
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
	m.InitConsole()
	m.MenuSelected = 0
}

func (m *Model) navigateTheme(direction int) {
	if len(m.AvailableThemes) == 0 {
		return
	}
	currentIdx := 0
	for i, t := range m.AvailableThemes {
		if m.Config != nil && t == m.Config.ThemeName {
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
	if m.Config != nil {
		oldTheme := m.Config.ThemeName
		m.Config.ThemeName = m.AvailableThemes[nextIdx]
		if m.Theme != nil {
			m.Theme.SetTintID(m.Config.ThemeName)
		}
		m.Styles = NewStyles(m.Theme)
		if err := m.SaveConfig(); err != nil {
			m.Config.ThemeName = oldTheme
			if m.Theme != nil {
				m.Theme.SetTintID(oldTheme)
			}
			m.Styles = NewStyles(m.Theme)
		}
	}
}
