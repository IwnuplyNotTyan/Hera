package generate

import (
	"github.com/charmbracelet/bubbles/key"

	"hera/i18n"
)

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Confirm key.Binding
	Shoot   key.Binding
	Ult     key.Binding
	Help    key.Binding
	Quit    key.Binding
	loc     i18n.Localizer
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Confirm, k.Shoot, k.Ult},
		{k.Help, k.Quit},
	}
}

func newKeyMap(loc i18n.Localizer) keyMap {
	return keyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp(loc.T("keys.up"), loc.T("help.moveUp")),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp(loc.T("keys.down"), loc.T("help.moveDown")),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp(loc.T("keys.left"), loc.T("help.moveLeft")),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp(loc.T("keys.right"), loc.T("help.moveRight")),
		),
		Confirm: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp(loc.T("keys.confirm"), loc.T("help.movePlayer")),
		),
		Shoot: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp(loc.T("keys.shoot"), loc.T("help.changeMode")),
		),
		Ult: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp(loc.T("keys.ult"), loc.T("help.secondAttack")),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp(loc.T("keys.help"), loc.T("help.toggleHelp")),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc", "ctrl+c"),
			key.WithHelp(loc.T("keys.quit"), loc.T("help.quit")),
		),
		loc: loc,
	}
}
