package ui

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

var borderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(table.DefaultStyles().Selected.GetForeground())
