package ui

import (
	"sort"
	"strings"

	"charm.land/bubbles/v2/table"
	"github.com/krabiworld/sshm/internal/config"
)

func (m *model) updateTable() {
	var rows []table.Row

	servers := m.config.GetAll()
	keys := make([]string, 0, len(servers))
	for server := range servers {
		keys = append(keys, server)
	}
	sort.Strings(keys)

	filter := m.searchInput.Value()

	for _, name := range keys {
		server := servers[name]

		if filter != "" {
			if !strings.Contains(name, filter) && !strings.Contains(server.Address, filter) {
				continue
			}
		}

		var identity string
		switch server.AuthType {
		case config.AuthPassword:
			identity = "Password"
		case config.AuthKey:
			identity = server.IdentityFile
		case config.AuthAgent:
			identity = "Agent"
		}

		rows = append(rows, table.Row{
			name,
			server.Address,
			server.Username,
			server.Port,
			identity,
		})
	}

	wasEmpty := len(m.table.Rows()) == 0

	m.table.SetRows(rows)

	if wasEmpty && len(rows) > 0 {
		m.table.GotoTop()
	}
}

func (m *model) recalculateTable() {
	var newColumns []table.Column
	for _, col := range columns {
		calculatedWidth := m.totalWidth / len(columns)

		newColumns = append(newColumns, table.Column{
			Title: col.Title,
			Width: calculatedWidth,
		})
	}

	m.table.SetColumns(newColumns)
	m.table.SetHeight(m.calculateHeight())
	m.table.SetWidth(m.calculateWidth())
}

func (m model) calculateWidth() int {
	return m.totalWidth
}

func (m model) calculateHeight() int {
	return m.totalHeight - 1
}
