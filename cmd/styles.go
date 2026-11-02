package cmd

import "charm.land/lipgloss/v2"

var (
	primaryColor = lipgloss.Color("#7B68EE")
	mutedColor   = lipgloss.Color("#8A8F98")
	textColor    = lipgloss.Color("#E6E6E6")
	successColor = lipgloss.Color("#4CDE8B")
	dangerColor  = lipgloss.Color("#FF5F87")

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor)

	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Width(12)

	valStyle = lipgloss.NewStyle().
			Foreground(textColor)

	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor)

	errStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(dangerColor)
)
