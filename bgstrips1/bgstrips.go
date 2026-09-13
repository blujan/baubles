// Package bgstrips1 gives stripy backgrounds
package bgstrips1

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"charm.land/bubbles/v2/help"

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
	keys       KeyMap
	help       help.Model
}

type (
	TickMsg time.Time
)

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 50)
	return TickMsg(time.Now())
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
		keys:       DefaultKeyMap(),
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
	for row := 0; row < m.height; row++ {
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
	return s.String()
}

func (m model) keyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		m.hPadding++
	case key.Matches(msg, m.keys.Down):
		m.hPadding = max(0, m.hPadding-1)
	case key.Matches(msg, m.keys.Left):
		m.wPadding = max(0, m.wPadding-1)
		return m.createSource(), nil
	case key.Matches(msg, m.keys.Right):
		m.wPadding++
		return m.createSource(), nil
	case key.Matches(msg, m.keys.DiagUp):
		m.wPadding++
		m.hPadding++
		return m.createSource(), nil
	case key.Matches(msg, m.keys.DiagDown):
		m.wPadding = max(0, m.wPadding-1)
		m.hPadding = max(0, m.hPadding-1)
		return m.createSource(), nil
	case key.Matches(msg, m.keys.Reset):
		m.wPadding = 1
		m.hPadding = 1
		return m.createSource(), nil
	}
	return m, nil
}

func (m model) KeyMap() help.KeyMap {
	return m.keys
}

func (m model) Status() string {
	return fmt.Sprintf("%d %d", m.wPadding, m.hPadding)
}

// BubbleTea API --------

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TickMsg:
		m.offset = (m.offset + 1) % len(m.source)
		return m, tick
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		return m.keyPress(msg)
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("")
	}
	return tea.NewView(m.createView())
}

// --------------
