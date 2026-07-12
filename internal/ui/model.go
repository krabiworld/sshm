package ui

import (
	"charm.land/bubbles/v2/table"
	"charm.land/bubbles/v2/textinput"
	"charm.land/huh/v2"
	"github.com/krabiworld/sshm/internal/config"
)

type modalType int

const (
	modalNone modalType = iota
	modalSearch
	modalCreate
	modalModify
	modalDelete
	modalSettings
)

var columns = []table.Column{
	{Title: "Name"},
	{Title: "Address"},
	{Title: "Username"},
	{Title: "Port"},
	{Title: "Identity"},
}

type model struct {
	config          config.Config
	configPath      string
	table           table.Model
	searchInput     textinput.Model
	form            *huh.Form
	activeModal     modalType
	totalWidth      int
	totalHeight     int
}

func NewModel(cfg config.Config, cfgPath string) model {
	m := model{
		config:     cfg,
		configPath: cfgPath,
		table: table.New(
			table.WithColumns(columns),
			table.WithFocused(true),
		),
		searchInput: textinput.New(),
	}

	m.updateTable()

	return m
}
