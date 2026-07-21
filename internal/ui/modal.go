package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/ui/forms"
)

func (m *model) wrapModal(msg tea.Msg, f func()) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted && m.form.GetBool(forms.Confirmed) {
		f()
	}

	if m.form.State == huh.StateCompleted || m.form.State == huh.StateAborted {
		m.activeModal = modalNone
	}

	return m, cmd
}

func (m *model) updateServer() {
	formName := m.form.GetString(forms.ServerName)
	formAddress := m.form.GetString(forms.ServerAddress)
	formUsername := m.form.GetString(forms.ServerUsername)
	formPort := m.form.GetString(forms.ServerPort)
	formAuthType := m.form.Get(forms.ServerAuthType).(config.AuthType)
	formIdentityFile := m.form.GetString(forms.ServerIdentityFile)
	formKnownHostsFile := m.form.GetString(forms.ServerKnownHostsFile)
	formPassword := m.form.GetString(forms.ServerPassword)

	server := config.Server{
		Address:        formAddress,
		Username:       formUsername,
		Port:           formPort,
		AuthType:       formAuthType,
		IdentityFile:   formIdentityFile,
		KnownHostsFile: formKnownHostsFile,
	}

	password := strings.TrimSpace(formPassword)
	if password != "" {
		encryptedPassword, err := m.crypto.Encrypt(password)
		if err != nil {
			panic(err)
		}
		server.Password = encryptedPassword
	}

	if m.activeModal == modalModify {
		currentName := m.table.SelectedRow()[0]
		if currentName != forms.ServerName {
			m.config.Delete(currentName)
		}
	}

	err := m.config.Set(formName, server)
	if err != nil {
		panic(err)
	}

	m.updateTable()
}

func (m *model) updateSettings() {
	formUsername := m.form.GetString(forms.ServerUsername)
	formPort := m.form.GetString(forms.ServerPort)
	formAuthType := m.form.Get(forms.ServerAuthType).(config.AuthType)
	formIdentityFile := m.form.GetString(forms.ServerIdentityFile)
	formKnownHostsFile := m.form.GetString(forms.ServerKnownHostsFile)

	if m.activeModal == modalModify {
		currentName := m.table.SelectedRow()[0]
		if currentName != forms.ServerName {
			m.config.Delete(currentName)
		}
	}

	defaults := m.config.GetDefaults()
	defaults.Username = formUsername
	defaults.Port = formPort
	defaults.AuthType = formAuthType
	defaults.IdentityFile = formIdentityFile
	defaults.KnownHostsFile = formKnownHostsFile

	if err := m.config.SetDefaults(defaults); err != nil {
		panic(err)
	}

	m.updateTable()
}
