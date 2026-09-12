package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/blujan/baubles/bgstrips1"
)

func main() {
	bg := lipgloss.Color("#1C1C1A")
	color1 := lipgloss.Color("#BB9AF7")
	color2 := lipgloss.Color("#2D364F")
	if _, err := tea.NewProgram(bgstrips1.New(bg, color1, color2, "░▒▓")).Run(); err != nil {
		// if _, err := tea.NewProgram(bgstrips1.New(bg, color1, color2, "▒▓", 5)).Run(); err != nil {

		fmt.Fprintln(os.Stderr, "Everything hurts: ", err)
		os.Exit(1)
	}
}
