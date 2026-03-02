package generate

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

const (
	GridW      = 14
	GridH      = 10
	wallCount  = 10
	waterCount = 10
	moveRange  = 4
	shootRange = 2
	MaxHP      = 3
)

type Point struct {
	X, Y int
}

type Player struct {
	X, Y  int
	HP    int
	Style lipgloss.Style
}

type Enemy struct {
	X, Y  int
	HP    int
	Style lipgloss.Style
}

type enemyTurnMsg struct {
	enemyIdx int
}

type Model struct {
	Players          []Player
	Enemys           []Enemy
	CurrentPlayer    int
	CurrentEnemy     int
	CursorX, CursorY int
	Walls            map[Point]bool
	Water            map[Point]bool
	ShootMode        bool
	Moved            bool
	Shot             bool
	keys             keyMap
	help             help.Model
	EnemyTurn        bool
	EnemyIdx         int
}

