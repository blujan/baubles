// Package bgstrips1 gives stripy backgrounds
package bgstrips1

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"github.com/blujan/baubles/cell"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	width      int
	height     int
	offset     int
	wPadding   int
	hPadding   int
	character  string
	background lipgloss.Style
	color1     lipgloss.Style
	color2     lipgloss.Style
	source     []string
	keys       keyMap
	help       help.Model
	styles     styles
}

type styles struct {
	Border      lipgloss.Style
	Title       lipgloss.Style
	Text        lipgloss.Style
	TextPressed lipgloss.Style
}

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	DiagUp   key.Binding
	DiagDown key.Binding
	Reset    key.Binding
	Quit     key.Binding
}

type (
	tickMsg time.Time
)

type unPressMsg struct {
	key tea.KeyPressMsg
}

var (
	borderColor      color.Color = lipgloss.Color("#8C5C46")
	borderTitleColor color.Color = lipgloss.Color("#ff9d65")
	bgColor          color.Color = lipgloss.Color("#1A1B26")
	textColor        color.Color = lipgloss.Color("#a8b1d6")
	textPressedColor color.Color = lipgloss.Color("#d63959")
)

var roundedBorder = lipgloss.Border{
	Top:          "─",
	Bottom:       "─",
	Left:         "│",
	Right:        "│",
	TopLeft:      "╭",
	TopRight:     "╮",
	BottomLeft:   "╰",
	BottomRight:  "╯",
	MiddleLeft:   "├",
	MiddleRight:  "┤",
	Middle:       "┼",
	MiddleTop:    "┬",
	MiddleBottom: "┴",
}

var keys = keyMap{
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
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},         // first column
		{k.DiagUp, k.DiagDown, k.Reset, k.Quit}, // second column
	}
}

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 50)
	return tickMsg(time.Now())
}

func unpress(key tea.KeyPressMsg) tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(_ time.Time) tea.Msg {
		return unPressMsg{key}
	})
}

func setupHelp() help.Model {
	mod := help.New()
	mod.ShowAll = true
	mod.Styles.FullKey = lipgloss.NewStyle().
		Foreground(borderTitleColor).
		Background(bgColor)
	mod.Styles.FullDesc = lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgColor)
	mod.Styles.FullSeparator = lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgColor)
	mod.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(textColor).
		Background(bgColor)
	return mod
}

func New(background, text1, text2 color.Color, char string) tea.Model {
	return model{
		offset:     0,
		wPadding:   1,
		hPadding:   1,
		character:  char,
		background: lipgloss.NewStyle().Background(background),
		color1:     lipgloss.NewStyle().Background(background).Foreground(text1),
		color2:     lipgloss.NewStyle().Background(background).Foreground(text2),
		keys:       keys,
		help:       setupHelp(),
		styles: styles{
			Border:      lipgloss.NewStyle().Background(bgColor).Foreground(borderColor),
			Title:       lipgloss.NewStyle().Background(bgColor).Foreground(borderTitleColor),
			Text:        lipgloss.NewStyle().Background(bgColor).Foreground(textColor),
			TextPressed: lipgloss.NewStyle().Background(bgColor).Foreground(textPressedColor),
		},
	}.createSource()
}

// Creates the string that's repeated to create the illusion of movement.
// Depends on m.wPadding. So, this needs to be re-run every time that value
// is updated. Go does not have operator overloading. So, we get to do that
// manually or create a separate function to inc/dec the value and then run
// this.
func (m model) createSource() tea.Model {
	runes := []rune(m.character)
	length := (m.wPadding + len(runes)) * 2
	m.source = make([]string, 0, length)
	for i := range runes {
		m.source = append(m.source, m.color1.Render(string(runes[i])))
	}
	for range m.wPadding {
		m.source = append(m.source, m.background.Render(" "))
	}
	for i := range runes {
		m.source = append(m.source, m.color2.Render(string(runes[i])))
	}
	for range m.wPadding {
		m.source = append(m.source, m.background.Render(" "))
	}
	return m
}

func (m model) createView() string {
	var s strings.Builder
	state := 0
	help := m.help.View(m.keys)
	for row := 0; row < (m.height - 6); row++ {
		switch state {
		case 0:
			for i := range m.width {
				index := (i + m.offset) % len(m.source)
				s.WriteString(m.source[index])
			}
		case m.hPadding + 1:
			for i := range m.width {
				index := (i - m.offset) % len(m.source)
				if index < 0 {
					index = len(m.source) + index
				}
				s.WriteString(m.source[index])
			}
		default:
			s.WriteString(m.background.Render(strings.Repeat(" ", m.width)))
		}
		state = (state + 1) % (2 + (m.hPadding * 2))
		s.WriteString("\n")

	}
	bottom := lipgloss.NewStyle().
		Width(m.width - 2).
		Background(bgColor).
		Align(lipgloss.Center).
		Render(help)
	border := cell.Cell{
		Border:          roundedBorder,
		TextTopLeft:     m.styles.Title.Render("Help"),
		TextBottomRight: m.styles.Text.Render(fmt.Sprintf("%d %d", m.wPadding, m.hPadding)),
		BorderColor:     borderColor,
		BgColor:         bgColor,
	}
	s.WriteString(border.Build(bottom))
	s.WriteString("\n")
	return s.String()
}

func (m model) keyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		m.hPadding++
		m.keys.Up.SetHelp(
			m.keys.Up.Help().Key,
			m.styles.TextPressed.Render(keys.Up.Help().Desc),
		)
		return m, unpress(msg)
	case key.Matches(msg, m.keys.Down):
		m.keys.Down.SetHelp(
			m.keys.Down.Help().Key,
			m.styles.TextPressed.Render(keys.Down.Help().Desc),
		)
		m.hPadding = max(0, m.hPadding-1)
		return m, unpress(msg)
	case key.Matches(msg, m.keys.Left):
		m.keys.Left.SetHelp(
			m.keys.Left.Help().Key,
			m.styles.TextPressed.Render(keys.Left.Help().Desc),
		)
		m.wPadding = max(0, m.wPadding-1)
		return m.createSource(), unpress(msg)
	case key.Matches(msg, m.keys.Right):
		m.keys.Right.SetHelp(
			m.keys.Right.Help().Key,
			m.styles.TextPressed.Render(keys.Right.Help().Desc),
		)
		m.wPadding++
		return m.createSource(), unpress(msg)
	case key.Matches(msg, m.keys.DiagUp):
		m.keys.DiagUp.SetHelp(
			m.keys.DiagUp.Help().Key,
			m.styles.TextPressed.Render(keys.DiagUp.Help().Desc),
		)
		m.wPadding++
		m.hPadding++
		return m.createSource(), unpress(msg)
	case key.Matches(msg, m.keys.DiagDown):
		m.keys.DiagDown.SetHelp(
			m.keys.DiagDown.Help().Key,
			m.styles.TextPressed.Render(keys.DiagDown.Help().Desc),
		)
		m.wPadding = max(0, m.wPadding-1)
		m.hPadding = max(0, m.hPadding-1)
		return m.createSource(), unpress(msg)
	case key.Matches(msg, m.keys.Reset):
		m.keys.Reset.SetHelp(
			m.keys.Reset.Help().Key,
			m.styles.TextPressed.Render(keys.Reset.Help().Desc),
		)
		m.wPadding = 1
		m.hPadding = 1
		return m.createSource(), unpress(msg)
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	}
	return m, nil
}

func (m model) keyUnpress(msg unPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg.key, m.keys.Up):
		m.keys.Up.SetHelp(
			m.keys.Up.Help().Key,
			// m.styles.Text.Render(keys.Up.Help().Desc),
			keys.Up.Help().Desc,
		)
	case key.Matches(msg.key, m.keys.Down):
		m.keys.Down.SetHelp(
			m.keys.Down.Help().Key,
			// m.styles.Text.Render(keys.Down.Help().Desc),
			keys.Down.Help().Desc,
		)
	case key.Matches(msg.key, m.keys.Left):
		m.keys.Left.SetHelp(
			m.keys.Left.Help().Key,
			// m.styles.Text.Render(keys.Left.Help().Desc),
			keys.Left.Help().Desc,
		)
	case key.Matches(msg.key, m.keys.Right):
		m.keys.Right.SetHelp(
			m.keys.Right.Help().Key,
			// m.styles.Text.Render(keys.Right.Help().Desc),
			keys.Right.Help().Desc,
		)
	case key.Matches(msg.key, m.keys.DiagUp):
		m.keys.DiagUp.SetHelp(
			m.keys.DiagUp.Help().Key,
			// m.styles.Text.Render(keys.DiagUp.Help().Desc),
			keys.DiagUp.Help().Desc,
		)
	case key.Matches(msg.key, m.keys.DiagDown):
		m.keys.DiagDown.SetHelp(
			m.keys.DiagDown.Help().Key,
			// m.styles.Text.Render(keys.DiagDown.Help().Desc),
			keys.DiagDown.Help().Desc,
		)
	case key.Matches(msg.key, m.keys.Reset):
		m.keys.Reset.SetHelp(
			m.keys.Reset.Help().Key,
			// m.styles.Text.Render(keys.Reset.Help().Desc),
			keys.Reset.Help().Desc,
		)
	case key.Matches(msg.key, m.keys.Quit):
		return m, tea.Quit
	}
	return m, nil
}

// BubbleTea API --------

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.offset = (m.offset + 1) % len(m.source)
		return m, tick
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		return m.keyPress(msg)
	case unPressMsg:
		return m.keyUnpress(msg)
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}
	var v tea.View
	v.AltScreen = true
	v.SetContent(m.createView())
	v.WindowTitle = "bgstrips1"
	return v
}

// --------------
