package generate

import (
	"github.com/charmbracelet/bubbles/help"

	"hera/i18n"
)

// NewModel creates a new Model with the given number of players and enemies, applying starter effects and theme.
func NewModel(playerCount, enemysCount int, playerEffects, enemyEffects []Effect, loc i18n.Localizer, theme ThemeRegistry, centerWindow string, enableBackground bool, themeName string) Model {
	styles := NewStyles(theme)
	if playerCount < 2 {
		playerCount = 2
	}
	if playerCount > 4 {
		playerCount = 4
	}
	if enemysCount < 0 {
		enemysCount = 0
	}
	if enemysCount > len(styles.EnemysStyles) {
		enemysCount = len(styles.EnemysStyles)
	}

	blocked := make(map[Point]bool)

	walls := GenerateTiles(GridW/2, GridH/2, wallCount, nil, nil)
	water := GenerateTiles(GridW/2, GridH/2, waterCount, walls, nil)

	starts := []Point{
		{X: 1, Y: 1},
		{X: GridW - 2, Y: GridH - 2},
		{X: GridW - 2, Y: 1},
		{X: 1, Y: GridH - 2},
	}

	players := make([]Player, playerCount)
	for i := range players {
		effs := make([]Effect, len(playerEffects))
		copy(effs, playerEffects)
		players[i] = Player{
			X:          starts[i].X,
			Y:          starts[i].Y,
			HP:         MaxHP,
			UltCharges: maxUltCharges,
			Effects:    effs,
			Style:      styles.PlayerStyles[i%len(styles.PlayerStyles)],
		}
	}

	for _, p := range players {
		blocked[Point{p.X, p.Y}] = true
	}
	for p := range walls {
		blocked[p] = true
	}
	for p := range water {
		blocked[p] = true
	}

	enemyStarts := GenerateTiles(GridW/2, GridH/2, enemysCount, blocked, nil)
	enemyPositions := make([]Point, 0, enemysCount)
	for p := range enemyStarts {
		enemyPositions = append(enemyPositions, p)
	}

	enemys := make([]Enemy, enemysCount)
	for i := range enemys {
		effs := make([]Effect, len(enemyEffects))
		copy(effs, enemyEffects)
		enemys[i] = Enemy{
			X:       enemyPositions[i].X,
			Y:       enemyPositions[i].Y,
			HP:      MaxHP,
			Effects: effs,
			Style:   styles.EnemysStyles[i],
		}
	}

	return Model{
		Theme:              theme,
		Config:             &Config{ThemeName: themeName, CenterWindow: centerWindow, Background: enableBackground},
		Styles:             styles,
		Screen:             ScreenMenu,
		EasterEgg:          loc.RandomEasterEgg(),
		Players:            players,
		Enemys:             enemys,
		CurrentPlayer:      0,
		CurrentEnemy:       0,
		CursorX:            players[0].X,
		CursorY:            players[0].Y,
		Walls:              walls,
		Water:              water,
		FireTiles:          make(map[Point]int),
		SmokeTiles:         make(map[Point]int),
		keys:               newKeyMap(loc),
		menuKeys:           newMenuKeyMap(loc),
		help:               help.New(),
		Localizer:          loc,
		startPlayers:       playerCount,
		startEnemies:       enemysCount,
		startPlayerEffects: playerEffects,
		startEnemyEffects:  enemyEffects,
	}
}

// SetAvailableThemes populates AvailableThemes from the theme registry if one is set.
func (m *Model) SetAvailableThemes() {
	if m.Theme != nil {
		m.AvailableThemes = m.Theme.TintIDs()
	}
}
