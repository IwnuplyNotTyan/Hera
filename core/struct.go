package generate

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	"github.com/lrstanley/bubblezone"

	"hera/i18n"
)

type EffectType string

const (
	GridW      = 14
	GridH      = 10
	wallCount  = 10
	waterCount = 10
	moveRange  = 4
	shootRange = 2
	MaxHP      = 3

	EffectWet   EffectType = "wet"
	EffectFire  EffectType = "fire"
	EffectSmoke EffectType = "smoke"

	maxUltCharges = 3
)

type Effect struct {
	Type     EffectType
	Duration int
}

type Point struct {
	X, Y int
}

type Player struct {
	X, Y       int
	HP         int
	Effects    []Effect
	UltCharges int
	Style      lipgloss.Style
}

type Enemy struct {
	X, Y    int
	HP      int
	Effects []Effect
	Style   lipgloss.Style
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
	FireTiles        map[Point]int
	SmokeTiles       map[Point]int
	ShootMode        bool
	UltMode          bool
	UltAxis          string
	Moved            bool
	Shot             bool
	keys             keyMap
	help             help.Model
	EnemyTurn        bool
	EnemyIdx         int
	Localizer        i18n.Localizer
	Z                *zone.Manager
	lastClickX       int
	lastClickY       int
}
