// Package cell provides a surrounding border + title/tray
package cell

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

type Cell struct {
	Border           lipgloss.Border
	TextTopLeft      string
	TextTopRight     string
	TextTopCenter    string
	TextBottomLeft   string
	TextBottomRight  string
	TextBottomCenter string
	BorderColor      color.Color
	BgColor          color.Color
	widths           cellWidths
}

type cellWidths struct {
	Top              int
	Bottom           int
	Left             int
	Right            int
	TopLeft          int
	TopRight         int
	BottomLeft       int
	BottomRight      int
	TextTopLeft      int
	TextTopRight     int
	TextTopCenter    int
	TextBottomLeft   int
	TextBottomRight  int
	TextBottomCenter int
}

type cellLine struct {
	TextLeft   string
	TextRight  string
	TextCenter string
	Left       string
	Mid        string
	Right      string
	Widths     cellLineWidth
}

type cellLineWidth struct {
	TextLeft   int
	TextRight  int
	TextCenter int
	Left       int
	Mid        int
	Right      int
}

type cellBuilder struct {
	Top     cellLine
	Bottom  cellLine
	Left    cellLine
	Right   cellLine
	FgColor color.Color
	BgColor color.Color
}

const textPadding = 2

func newCellBuilder(cell Cell) (builder cellBuilder) {
	builder.Top = cellLine{
		TextLeft:   cell.TextTopLeft,
		TextRight:  cell.TextTopRight,
		TextCenter: cell.TextTopCenter,
		Left:       cell.Border.TopLeft,
		Right:      cell.Border.TopRight,
		Mid:        cell.Border.Top,
	}
	builder.Top.Widths = builderWidths(builder.Top)
	builder.Bottom = cellLine{
		TextLeft:   cell.TextBottomLeft,
		TextRight:  cell.TextBottomRight,
		TextCenter: cell.TextBottomCenter,
		Left:       cell.Border.BottomLeft,
		Right:      cell.Border.BottomRight,
		Mid:        cell.Border.Bottom,
	}
	builder.Bottom.Widths = builderWidths(builder.Bottom)
	builder.Left = cellLine{
		Left: cell.Border.Left,
	}
	builder.Left.Widths = builderWidths(builder.Left)
	builder.Right = cellLine{
		Right: cell.Border.Right,
	}
	builder.Right.Widths = builderWidths(builder.Right)
	builder.FgColor = cell.BorderColor
	builder.BgColor = cell.BgColor
	return
}

func builderWidths(builder cellLine) cellLineWidth {
	return cellLineWidth{
		TextLeft:   ansi.StringWidth(builder.TextLeft),
		TextRight:  ansi.StringWidth(builder.TextRight),
		TextCenter: ansi.StringWidth(builder.TextCenter),
		Left:       ansi.StringWidth(builder.Left),
		Mid:        ansi.StringWidth(builder.Mid),
		Right:      ansi.StringWidth(builder.Right),
	}
}

// Split a string into lines, additionally returning the size of the widest
// line.
// Mostly ripped from charmbracelet.lipgloss get.go.
// AI agents seem roll their own max/min functions rather than
// use the built-ins. That is fixed here since this
// is a human-slop-only zone.
func getLines(s string) (lines []string, widest int) {
	s = strings.ReplaceAll(s, "\t", "    ")
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		widest = max(widest, ansi.StringWidth(l))
	}
	return
}

// Go does not have constructors. So we have to manually call this before
// we can use our data structure or else the data will not be coherent.
func (cell Cell) setWidths() Cell {
	cell.widths = cellWidths{
		Top:              ansi.StringWidth(cell.Border.Top),
		Bottom:           ansi.StringWidth(cell.Border.Bottom),
		Left:             ansi.StringWidth(cell.Border.Left),
		Right:            ansi.StringWidth(cell.Border.Right),
		TopLeft:          ansi.StringWidth(cell.Border.TopLeft),
		TopRight:         ansi.StringWidth(cell.Border.TopRight),
		BottomLeft:       ansi.StringWidth(cell.Border.BottomLeft),
		BottomRight:      ansi.StringWidth(cell.Border.BottomRight),
		TextTopLeft:      ansi.StringWidth(cell.TextTopLeft),
		TextTopRight:     ansi.StringWidth(cell.TextTopRight),
		TextTopCenter:    ansi.StringWidth(cell.TextTopCenter),
		TextBottomLeft:   ansi.StringWidth(cell.TextBottomLeft),
		TextBottomRight:  ansi.StringWidth(cell.TextBottomRight),
		TextBottomCenter: ansi.StringWidth(cell.TextBottomCenter),
	}
	return cell
}

func (cell cellLine) build(target int, FgColor, BgColor color.Color) string {
	var out strings.Builder
	cornerStyle := lipgloss.NewStyle().Foreground(lipgloss.Lighten(FgColor, 0.05)).Background(BgColor)
	topStyle := lipgloss.NewStyle().Foreground(FgColor).Background(BgColor)
	width := 0
	centerStart := (target / 2) - (cell.Widths.TextCenter / 2) - 1
	rightStart := target - textPadding - cell.Widths.Mid - cell.Widths.TextRight - cell.Widths.Right
	done := false

	out.WriteString(cornerStyle.Render(cell.Left))
	out.WriteString(topStyle.Render(cell.Mid))
	width += (cell.Widths.Left + cell.Widths.Mid)
	if (width + cell.Widths.TextLeft + textPadding) < (target - cell.Widths.Right) {
		if cell.Widths.TextLeft > 0 {
			out.WriteString(topStyle.Render(" "))
			out.WriteString(cell.TextLeft)
			out.WriteString(topStyle.Render(" "))
			width += textPadding + cell.Widths.TextLeft
		}
	} else {
		done = true
	}
	// Center Text
	if width < centerStart {
		if !done && cell.Widths.TextCenter > 0 {
			count := (centerStart - width) / cell.Widths.Mid
			rem := (centerStart - width) % cell.Widths.Mid
			out.WriteString(topStyle.Render(strings.Repeat(cell.Mid, count)))
			width += (count * cell.Widths.Mid)
			if rem != 0 {
				runes := []rune(cell.Mid)
				out.WriteString(topStyle.Render(string(runes[0:rem])))
				width += rem
			}
			out.WriteString(topStyle.Render(" "))
			out.WriteString(cell.TextCenter)
			out.WriteString(topStyle.Render(" "))
			width += textPadding + cell.Widths.TextCenter
		}
	} else {
		done = true
	}
	// Left Text
	if width < rightStart {
		if !done && cell.Widths.TextRight > 0 {
			count := (rightStart - width) / cell.Widths.Mid
			rem := (rightStart - width) % cell.Widths.Mid
			out.WriteString(topStyle.Render(strings.Repeat(cell.Mid, count)))
			width += (count * cell.Widths.Mid)
			if rem != 0 {
				runes := []rune(cell.Mid)
				out.WriteString(topStyle.Render(string(runes[0:rem])))
				width += rem
			}
			out.WriteString(topStyle.Render(" "))
			out.WriteString(cell.TextRight)
			out.WriteString(topStyle.Render(" "))
			out.WriteString(topStyle.Render(cell.Mid))
			width += textPadding + cell.Widths.Mid + cell.Widths.TextRight
		}
	}
	// End, Right Corner
	if width < (target - cell.Widths.Right) {
		count := (target - cell.Widths.Right - width) / cell.Widths.Mid
		rem := (target - cell.Widths.Right - width) % cell.Widths.Mid
		out.WriteString(topStyle.Render(strings.Repeat(cell.Mid, count)))
		width += (target - cell.Widths.Right - width) * cell.Widths.Mid
		if rem != 0 {
			runes := []rune(cell.Mid)
			out.WriteString(topStyle.Render(string(runes[0:rem])))
			width += rem
		}
	}

	out.WriteString(cornerStyle.Render(cell.Right))

	return out.String()
}

func (builder cellBuilder) String(contents string) string {
	var out strings.Builder
	lines, textWidth := getLines(contents)
	width := builder.Left.Widths.Left + textWidth + builder.Right.Widths.Right
	out.WriteString(builder.Top.build(width, builder.FgColor, builder.BgColor))
	out.WriteString("\n")

	lineStyle := lipgloss.NewStyle().Width(textWidth)
	bgStyle := lipgloss.NewStyle().Foreground(builder.FgColor).Background(builder.BgColor)

	for _, line := range lines {
		out.WriteString(bgStyle.Render(builder.Left.Left))
		out.WriteString(lineStyle.Render(line))
		out.WriteString(bgStyle.Render(builder.Right.Right))
		out.WriteString("\n")
	}
	out.WriteString(builder.Bottom.build(width, builder.FgColor, builder.BgColor))
	out.WriteString("\n")
	return out.String()
}

func (cell Cell) Build(contents string) string {
	builder := newCellBuilder(cell)
	return builder.String(contents)
}
