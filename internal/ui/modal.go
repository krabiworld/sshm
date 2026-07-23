package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/ui/forms"
)

func (m *model) wrapModal(msg tea.Msg, f func() tea.Cmd) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted && m.form.GetBool(forms.Confirmed) {
		if cmd := f(); cmd != nil {
			return m, cmd
		}
	}

	if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
		m.activeModal = modalNone
	}

	return m, cmd
}

func (m *model) updateServer() tea.Cmd {
	formName := m.form.GetString(forms.ServerName)
	formAddress := m.form.GetString(forms.ServerAddress)
	formUsername := m.form.GetString(forms.ServerUsername)
	formPort := m.form.GetString(forms.ServerPort)
	formAuthType := m.form.Get(forms.ServerAuthType).(config.AuthType)
	formIdentityFile := m.form.GetString(forms.ServerIdentityFile)
	formPassword := m.form.GetString(forms.ServerPassword)

	server := config.Server{
		Address:      formAddress,
		Username:     formUsername,
		Port:         formPort,
		AuthType:     formAuthType,
		IdentityFile: formIdentityFile,
	}

	password := strings.TrimSpace(formPassword)
	if password != "" {
		encryptedPassword, err := m.crypto.Encrypt(password)
		if err != nil {
			return errCmd(err)
		}
		server.Password = encryptedPassword
	}

	if m.activeModal == modalModify {
		currentName := m.table.SelectedRow()[0]
		if currentName != formName {
			m.config.Delete(currentName)
		}

		currentPassword := m.config.Get(currentName).Password
		if password == "" && currentPassword != "" {
			server.Password = currentPassword
		}
	}

	err := m.config.Set(formName, server)
	if err != nil {
		return errCmd(err)
	}

	m.updateTable()

	return nil
}
