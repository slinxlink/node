package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor = lipgloss.Color("#FFB6B6")
	successColor = lipgloss.Color("#5ADC5A")
	errorColor   = lipgloss.Color("#DC2020")
)

func box(content string, title ...string) string {
	t := "SLINX"
	if len(title) > 0 && title[0] != "" {
		t = title[0]
	}
	t = " " + t + " "

	b := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render(content)

	topLine := strings.Split(b, "\n")[0]
	totalDashes := strings.Count(topLine, "─")
	boldTitle := lipgloss.NewStyle().Bold(true).Foreground(primaryColor).Render(t)
	newTop := "╭" + boldTitle + strings.Repeat("─", totalDashes-lipgloss.Width(t)) + "╮"
	return strings.Replace(b, topLine, newTop, 1)
}

func renderStatus(left, right string, success bool) string {
	var rightStyle lipgloss.Style
	if success {
		rightStyle = lipgloss.NewStyle().Foreground(successColor)
	} else {
		rightStyle = lipgloss.NewStyle().Foreground(errorColor)
	}
	content := left + " - " + rightStyle.Render(right)
	return box(content)
}

func Status(left, right string, success bool) {
	fmt.Println(renderStatus(left, right, success))
}

func renderInfo(title string, rows ...[]string) string {
	maxWidth := 0
	for _, row := range rows {
		w := lipgloss.Width(row[0])
		if w > maxWidth {
			maxWidth = w
		}
	}

	lines := ""
	for i, row := range rows {
		left := row[0]
		current := lipgloss.Width(left)
		if current < maxWidth {
			left = left + strings.Repeat(" ", maxWidth-current)
		}
		line := left + " - " + row[1]
		if i < len(rows)-1 {
			line += "\n"
		}
		lines += line
	}
	return box(lines, title)
}

func Info(title string, rows ...[]string) {
	fmt.Println(renderInfo(title, rows...))
}
