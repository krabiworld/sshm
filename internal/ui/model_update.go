package ui

import (
	"errors"
	"fmt"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/crypto"
	"github.com/krabiworld/sshm/internal/ui/forms"
	"github.com/krabiworld/sshm/internal/utils"
	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/ssh"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height

		m.recalculateTable()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.crypto == (crypto.Cipher{}) {
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		if m.form.State == huh.StateCompleted {
			masterPassword := m.form.GetString(forms.MasterPassword)
			m.crypto = crypto.NewCipher(masterPassword)

			if m.form.GetBool(forms.Confirmed) {
				err := keyring.Set(keyringService, keyringUsername, masterPassword)
				if err != nil {
					panic(err)
				}
			}
		}
		return m, cmd
	}

	// Global
	switch msg := msg.(type) {
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
		return m.wrapModal(msg, m.updateServer)
	case modalDelete:
		return m.wrapModal(msg, func() tea.Cmd {
			if err := m.config.Delete(m.getCurrentServer()); err != nil {
				return errCmd(err)
			}

			m.updateTable()

			return nil
		})
	case modalHostKeyRequired:
		return m.wrapModal(msg, func() tea.Cmd {
			if hostKeyErr, ok := errors.AsType[*hostKeyRequiredError](m.error); ok {
				err := utils.AddHostKey(hostKeyErr.hostname, hostKeyErr.key)
				if err != nil {
					return errCmd(err)
				}
				m.activeModal = modalConnecting
				return tea.Batch(m.dialSsh(m.getCurrentServer()), m.spinner.Tick)
			}

			return nil
		})
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
		}
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}
