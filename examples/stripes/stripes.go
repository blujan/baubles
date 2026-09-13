package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/blujan/baubles/bgstrips1"
	"github.com/blujan/baubles/cell"
	"github.com/charmbracelet/x/ansi"
)

//
// Structs
//

type model struct {
	background []Background
	bgIndex    int
	keys       keymap
	styles     styles
	help       help.Model
	width      int
	height     int
	helpHeight int
}

type styles struct {
	Border      lipgloss.Style
	Title       lipgloss.Style
	Text        lipgloss.Style
	TextPressed lipgloss.Style
}

//
// Interfaces
//

type Background interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
	KeyMap() help.KeyMap
	Status() string
}

//
// Constants
//

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

//
// Events
//

type UnpressMsg struct {
	key  tea.KeyPressMsg
	desc string
}

func unpress(key tea.KeyPressMsg, desc string) tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(_ time.Time) tea.Msg {
		return UnpressMsg{key, desc}
	})
}

//
// Keymaps
//

type keymap struct {
	Quit key.Binding
	// short func() []key.Binding
	// full  func() [][]key.Binding
	short []key.Binding
	full  [][]key.Binding
}

var keys = keymap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (k keymap) ShortHelp() []key.Binding {
	return k.short
}

func (k keymap) FullHelp() [][]key.Binding {
	return k.full
}

//
// Helper functions
//

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

func Flatten[T any](lists [][]T) []T {
	var res []T
	for _, list := range lists {
		res = append(res, list...)
	}
	return res
}

//
// Model
//

func New() tea.Model {
	m := model{
		bgIndex: 0,
		keys:    keys,
		help:    setupHelp(),
		styles: styles{
			Border:      lipgloss.NewStyle().Background(bgColor).Foreground(borderColor),
			Title:       lipgloss.NewStyle().Background(bgColor).Foreground(borderTitleColor),
			Text:        lipgloss.NewStyle().Background(bgColor).Foreground(textColor),
			TextPressed: lipgloss.NewStyle().Background(bgColor).Foreground(textPressedColor),
		},
		helpHeight: 4,
	}
	m.background = append(m.background,
		bgstrips1.New(bgColor,
			lipgloss.Color("#BB9AF7"),
			lipgloss.Color("#2D364F"),
			"░▒▓").(Background))
	m.GenerateHelp()
	return m
}

func (m *model) GenerateHelp() {
	bgKeyMap := m.background[m.bgIndex].KeyMap()
	short := bgKeyMap.ShortHelp()
	long := Flatten(bgKeyMap.FullHelp())
	long = append(long, m.keys.Quit)

	var full [][]key.Binding
	full = append(full, make([]key.Binding, 0))
	index := 0
	count := 0
	for _, value := range long {
		full[index] = append(full[index], value)
		count++
		if count == m.helpHeight {
			count = 0
			full = append(full, make([]key.Binding, 0))
			index++
		}
	}

	m.keys.short = short
	m.keys.full = full
}

func (m model) showKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	for rindex, row := range m.keys.full {
		for index, item := range row {
			if key.Matches(msg, item) {
				desc := item.Help().Desc
				m.keys.full[rindex][index].SetHelp(
					item.Help().Key,
					m.styles.TextPressed.Render(desc),
				)
				return unpress(msg, desc)
			}
		}
	}
	return nil
}

func (m model) showKeyUnpress(msg UnpressMsg) {
	for rindex, row := range m.keys.full {
		for index, item := range row {
			if key.Matches(msg.key, item) {
				m.keys.full[rindex][index].SetHelp(
					item.Help().Key,
					ansi.Strip(msg.desc),
				)
				return
			}
		}
	}
}

// ---------------------

func (m model) Init() tea.Cmd {
	return m.background[m.bgIndex].Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		bg, bgCmd := m.background[m.bgIndex].Update(
			tea.WindowSizeMsg{
				Width:  msg.Width,
				Height: msg.Height - m.helpHeight - 2,
			})
		m.background[m.bgIndex] = bg.(Background)
		return m, tea.Batch(cmd, bgCmd)
	case tea.KeyPressMsg:
		cmd = m.showKeyPress(msg)
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	case UnpressMsg:
		m.showKeyUnpress(msg)
		return m, nil
	}
	bg, bgCmd := m.background[m.bgIndex].Update(msg)
	m.background[m.bgIndex] = bg.(Background)
	return m, tea.Batch(cmd, bgCmd)
}

func (m model) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.WindowTitle = "Stripes Stripes Stripes"

	bg := m.background[m.bgIndex].View()
	help := m.help.View(m.keys)
	bottom := lipgloss.NewStyle().
		Width(m.width - 2).
		Height(m.helpHeight).
		Background(bgColor).
		Align(lipgloss.Center).
		Render(help)
	border := cell.Cell{
		Border:          roundedBorder,
		TextTopLeft:     m.styles.Title.Render("Help"),
		TextBottomRight: m.styles.Text.Render(m.background[m.bgIndex].Status()),
		BorderColor:     borderColor,
		BgColor:         bgColor,
	}

	var s strings.Builder
	s.WriteString(bg.Content)
	s.WriteString(border.Build(bottom))
	s.WriteString("\n")
	v.Content = s.String()
	return v
}

func main() {
	if _, err := tea.NewProgram(New()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Oh no", err)
		os.Exit(1)
	}
}
