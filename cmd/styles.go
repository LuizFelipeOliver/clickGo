package cmd

import "charm.land/lipgloss/v2"

var (
	primary = lipgloss.Color("#7B68EE")
	muted   = lipgloss.Color("#8A8F98")
	text    = lipgloss.Color("#E6E6E6")
	success = lipgloss.Color("#4CDE8B")
	warn    = lipgloss.Color("#FFC15E")
	danger  = lipgloss.Color("#FF5F87")

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primary)

	keyStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primary).
		Width(12)

	valStyle = lipgloss.NewStyle().
		Foreground(text)

	mutedStyle = lipgloss.NewStyle().
		Foreground(muted)

	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(success)

	warnStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(warn)

	errStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(danger)
)
