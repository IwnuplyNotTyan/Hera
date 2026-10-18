package generate

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"hera/i18n"
)

type Effect struct {
	Type     EffectType
	Duration int
}

// Wall represents a destructible stone on the grid with hit points.
type Wall struct {
	HP int
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
	AttackType AttackType
}

type Enemy struct {
	X, Y       int
	HP         int
	Effects    []Effect
	Style      lipgloss.Style
	AttackType AttackType
}

type enemyTurnMsg struct {
	enemyIdx int
}

type triggerTickMsg struct{}

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
	Walls                    map[Point]Wall
	Water                    map[Point]bool
	FireTiles                map[Point]int
	SmokeTiles               map[Point]int
	ShowEffectIdx            int
	ShootMode                bool
	PushStrikeMode           bool
	RamMode                  bool
	MeleePushMode            bool
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
	CenterRow, CenterCol     int

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
	Seed                 int64
	SeedLocked           bool
	SeedConfirmActive    bool
	SeedConfirmChoice    int
	SeedPromptInput      textinput.Model
	DebugMode            bool

	ConsoleMode       bool
	ConsoleInput      textinput.Model
	ConsoleOutput     []string
	ConsoleHistory    []string
	ConsoleHistoryIdx int
}
