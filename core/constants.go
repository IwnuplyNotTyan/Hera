package generate

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

type ElementType int

const (
	ElementGridCell ElementType = iota
	ElementMenuItem
	ElementSettingsItem
	ElementThemeItem
	ElementProfileConfirm
	ElementCenterItem
)

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenSettings
	ScreenThemeSelect
	ScreenCenterSelect
	ScreenGame
	ScreenProfiles
	ScreenProfileCreate
	ScreenGameOver
	ScreenSeedPrompt
)

const (
	TriggerNone   = ""
	TriggerDamage = "damage"
)
