package generate

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"

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

type triggerTickMsg struct{}

type ElementType int

const (
	ElementGridCell ElementType = iota
	ElementMenuItem
	ElementSettingsItem
	ElementThemeItem
	ElementProfileConfirm
)

type Element struct {
	Type          ElementType
	X, Y          int
	Width, Height int
	ID            string
	Index         int
}

type Profile struct {
	Name  string
	Score int
}

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenSettings
	ScreenThemeSelect
	ScreenGame
	ScreenProfiles
	ScreenProfileCreate
)

type Model struct {
	Theme                    ThemeRegistry
	Styles                   Styles
	Config                   *Config
	TerminalWidth            int
	TerminalHeight           int
	Screen                   Screen
	MenuSelected             int
	EasterEgg                string
	ThemeSearch              bool
	ThemeQuery               string
	LastSearchQuery          string
	AvailableThemes          []string
	Players                  []Player
	Enemys                   []Enemy
	CurrentPlayer            int
	CurrentEnemy             int
	CursorX, CursorY         int
	Walls                    map[Point]bool
	Water                    map[Point]bool
	FireTiles                map[Point]int
	SmokeTiles               map[Point]int
	ShowEffectIdx            int
	ShootMode                bool
	UltMode                  bool
	UltAxis                  string
	Moved                    bool
	Shot                     bool
	keys                     keyMap
	help                     help.Model
	menuKeys                 menuKeyMap
	EnemyTurn                bool
	EnemyIdx                 int
	Localizer                i18n.Localizer
	layoutElements           []Element
	gridOffsetX, gridOffsetY int

	BannerText   string
	BoxTrigger   string
	TriggerTimer int

	startPlayers       int
	startEnemies       int
	startPlayerEffects []Effect
	startEnemyEffects  []Effect

	Profiles             [3]*Profile
	ProfileSlot          int
	ProfileLetters       [3]rune
	ProfileCursor        int
	Font                 *FLFFont
	ProfileDeleteConfirm bool
	ProfileConfirmChoice int
	HoveredConfirm       string
	CurrentScore         int
	ActiveSlot           int
}

const (
	TriggerNone   = ""
	TriggerDamage = "damage"
)
