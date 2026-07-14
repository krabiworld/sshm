package ui

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

var paddingStyle = lipgloss.NewStyle().
	Padding(0, 1)

var borderStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(table.DefaultStyles().Selected.GetForeground())
