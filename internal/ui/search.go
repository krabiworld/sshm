package ui

import tea "charm.land/bubbletea/v2"

func (m model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc", "enter":
			m.activeModal = modalNone
			m.searchInput.Blur()
			m.table.Focus()
			return m, nil
		}
	}

	m.searchInput, cmd = m.searchInput.Update(msg)
	m.updateTable()
	return m, cmd
}
