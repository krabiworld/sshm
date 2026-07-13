package ui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/ui/forms"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		if keyMsg.String() == "ctrl+c" {
			return m, tea.Quit
		} else if keyMsg.String() == "esc" && m.activeModal != modalNone {
			m.activeModal = modalNone
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height

		m.recalculateTable()
	}

	switch m.activeModal {
	case modalSearch:
		return m.updateSearch(msg)
	case modalCreate, modalModify:
		return m.updateServer(msg)
	case modalDelete:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted && m.form.GetBool(forms.DeleteConfirmed) {
			err := m.config.Delete(m.getCurrentServer(), m.configPath)
			if err != nil {
				panic(err)
			}
			m.updateTable()
		}
		if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
			m.activeModal = modalNone
		}
		return m, cmd
	case modalSettings:
		return m.updateSettings(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)
		case "enter":
			return m, m.connectSsh(m.getCurrentServer())
		case "ctrl+f":
			m.activeModal = modalSearch
			m.searchInput.Focus()
			return m, textinput.Blink
		case "ctrl+a":
			m.activeModal = modalCreate
			m.form = forms.NewServer(m.config, "")
			return m, m.form.Init()
		case "ctrl+m":
			m.activeModal = modalModify
			m.form = forms.NewServer(m.config, m.getCurrentServer())
			return m, m.form.Init()
		case "ctrl+d":
			m.activeModal = modalDelete
			m.form = forms.NewDelete(m.getCurrentServer())
			return m, m.form.Init()
		case "ctrl+i":
			return m, m.copyId(m.getCurrentServer())
		case "ctrl+x":
			m.activeModal = modalSettings
			m.form = forms.NewSettings(m.config)
			return m, m.form.Init()
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
