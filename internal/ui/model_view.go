package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	footer := " ^F Search | ^A Add | ^M Modify | ^D Delete | ^I Copy ID | ^X Settings"

	if m.activeModal == modalSearch {
		footer = lipgloss.JoinHorizontal(lipgloss.Left, " Search", m.searchInput.View())
	}

	layers := []*lipgloss.Layer{
		lipgloss.NewLayer(lipgloss.JoinVertical(
			lipgloss.Left,
			m.table.View(),
			footer,
		)),
	}

	if m.activeModal != modalNone && m.activeModal != modalSearch && m.activeModal != modalError {
		layers = m.appendModal(layers, m.form.View())
	}
	if m.activeModal == modalError {
		textWidth := m.totalWidth - 40

		errorStyle := lipgloss.NewStyle().Width(textWidth)
		errorContent := errorStyle.Render("Error:\n" + m.errorMessage)

		layers = m.appendModal(layers, errorContent)
	}

	view := tea.NewView(lipgloss.NewCompositor(layers...).Render())
	view.AltScreen = true
	return view
}

func (m model) appendModal(layers []*lipgloss.Layer, view string) []*lipgloss.Layer {
	safeMaxWidth := m.totalWidth - 4
	safeMaxHeight := m.totalHeight - 4

	modalRendered := borderStyle.
		MaxWidth(safeMaxWidth).
		MaxHeight(safeMaxHeight).
		Render(view)

	modalWidth := lipgloss.Width(modalRendered)
	modalHeight := lipgloss.Height(modalRendered)

	modalX := (m.totalWidth / 2) - (modalWidth / 2)
	modalY := (m.totalHeight / 2) - (modalHeight / 2)

	return append(layers, lipgloss.NewLayer(modalRendered).X(modalX).Y(modalY))
}
