package ui

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/security"
	"github.com/krabiworld/sshm/internal/ui/forms"
	"github.com/krabiworld/sshm/internal/utils"
	"golang.org/x/crypto/ssh"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height

		m.recalculateTable()
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case *sshConnectedMsg:
		m.activeModal = modalNone
		return m, m.runSshSession(msg)
	case errMsg:
		if hostKeyErr, ok := errors.AsType[*hostKeyRequiredError](msg.err); ok {
			m.activeModal = modalHostKeyRequired
			m.error = hostKeyErr
			m.form = forms.NewConfirm(fmt.Sprintf(
				"%s key fingerprint is: %s\n"+
					"This key is not known by any other names.\n"+
					"Are you sure you want to continue connecting?",
				hostKeyErr.key.Type(), ssh.FingerprintSHA256(hostKeyErr.key)))
			return m, m.form.Init()
		}
		m.error = msg.err
		m.activeModal = modalError
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.activeModal == modalConnecting {
				return m, nil
			}

			if m.activeModal != modalNone {
				m.activeModal = modalNone
				return m, nil
			}
		}
	}

	// Modal
	switch m.activeModal {
	case modalConnecting:
		return m, nil
	case modalSearch:
		return m.updateSearch(msg)
	case modalCreate, modalModify:
		return m.updateServer(msg)
	case modalDelete:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted && m.form.GetBool(forms.FormConfirmed) {
			serverName := m.getCurrentServer()
			if err := m.config.Delete(serverName); err != nil {
				return m, errCmd(err)
			}
			if err := security.DeletePassword(serverName); err != nil {
				return m, errCmd(err)
			}
			m.updateTable()
		}
		if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
			m.activeModal = modalNone
		}
		return m, cmd
	case modalSettings:
		return m.updateSettings(msg)
	case modalHostKeyRequired:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted && m.form.GetBool(forms.FormConfirmed) {
			if hostKeyErr, ok := errors.AsType[*hostKeyRequiredError](m.error); ok {
				err := utils.AddHostKey(hostKeyErr.path, hostKeyErr.hostname, hostKeyErr.key)
				if err != nil {
					return m, errCmd(err)
				}
				m.activeModal = modalConnecting
				return m, tea.Batch(m.dialSsh(m.getCurrentServer()), m.spinner.Tick)
			}
		}
		if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
			m.activeModal = modalNone
		}
		return m, cmd
	case modalError:
		if keyMsg, ok := msg.(tea.KeyPressMsg); ok && keyMsg.String() == "enter" {
			m.activeModal = modalNone
		}
		return m, nil
	}

	// Table
	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch keyMsg.String() {
		case "enter":
			m.activeModal = modalConnecting
			return m, tea.Batch(m.dialSsh(m.getCurrentServer()), m.spinner.Tick)
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
			m.form = forms.NewConfirm(fmt.Sprintf("Are you sure you want to delete %s?", m.getCurrentServer()))
			return m, m.form.Init()
		case "ctrl+i":
			return m, m.copyId(m.getCurrentServer())
		case "ctrl+x":
			m.activeModal = modalSettings
			m.form = forms.NewSettings(m.config)
			return m, m.form.Init()
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}
