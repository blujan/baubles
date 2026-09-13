package bgstrips1

import (
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	DiagUp   key.Binding
	DiagDown key.Binding
	Reset    key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "8"),
			key.WithHelp("↑/8", "+Vertical Distance"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "2"),
			key.WithHelp("↓/2", "-Vertical Distance"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "4"),
			key.WithHelp("←/4", "-Horizontal Distance"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "6"),
			key.WithHelp("→/6", "+Horizontal Distance"),
		),
		DiagUp: key.NewBinding(
			key.WithKeys("9", "7"),
			key.WithHelp("9/7", "+Vertical, +Horizontal Distance"),
		),
		DiagDown: key.NewBinding(
			key.WithKeys("1", "3"),
			key.WithHelp("1/3", "-Vertical, -Horizontal Distance"),
		),
		Reset: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "Reset"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // first column
		{k.DiagUp, k.DiagDown, k.Reset}, // second column
	}
}
