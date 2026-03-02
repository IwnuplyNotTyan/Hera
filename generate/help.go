package generate

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Confirm key.Binding
	Shoot   key.Binding
	Ult	key.Binding
	Help    key.Binding
	Quit    key.Binding
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

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/K", "Move Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/J", "Move Down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/H", "Move Left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/L", "Move Right"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("X", "Move Player"),
	),
	Shoot: key.NewBinding(
		key.WithKeys("z"),
		key.WithHelp("Z", "Change Mode"),
	),
	Ult: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("C", "Second Attack"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle Help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
}

