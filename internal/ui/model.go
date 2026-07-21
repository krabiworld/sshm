package ui

import (
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
	"github.com/krabiworld/sshm/internal/crypto"
	"github.com/krabiworld/sshm/internal/ui/forms"
	"github.com/zalando/go-keyring"
)

type modalType int

const (
	modalNone modalType = iota
	modalSearch
	modalCreate
	modalModify
	modalDelete
	modalSettings
	modalConnecting
	modalHostKeyRequired
	modalError

	keyringService  = "SSHM Safe Storage"
	keyringUsername = "SSHM"
)

var columns = []table.Column{
	{Title: "Name"},
	{Title: "Address"},
	{Title: "Username"},
	{Title: "Port"},
	{Title: "Identity"},
}

type model struct {
	config      *config.Config
	table       table.Model
	searchInput textinput.Model
	form        *huh.Form
	spinner     spinner.Model
	crypto      crypto.Cipher
	activeModal modalType
	totalWidth  int
	totalHeight int
	error       error
}

func NewModel(cfg *config.Config) model {
	m := model{
		config: cfg,
		table: table.New(
			table.WithColumns(columns),
			table.WithFocused(true),
		),
		searchInput: textinput.New(),
		form:        forms.NewMasterPassword(),
		spinner:     spinner.New(),
	}

	masterPassword, _ := keyring.Get(keyringService, keyringUsername)
	if masterPassword != "" {
		m.crypto = crypto.NewCipher(masterPassword)
	}

	m.spinner.Spinner = spinner.Dot
	m.updateTable()

	return m
}
