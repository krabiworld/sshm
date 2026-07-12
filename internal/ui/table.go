package ui

import (
	"sort"
	"strings"

	"charm.land/bubbles/v2/table"
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
		if server.HasPassword {
			identity = "Password"
		} else {
			identity = m.appendDefaultTag(server.IdentityFile, m.config.Defaults.IdentityFile)
		}

		rows = append(rows, table.Row{
			name,
			server.Address,
			m.appendDefaultTag(server.Username, m.config.Defaults.Username),
			m.appendDefaultTag(server.Port, m.config.Defaults.Port),
			identity,
		})
	}

	m.table.SetRows(rows)
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

func (m model) appendDefaultTag(val, def string) string {
	if val == def {
		return val + " (d)"
	}
	return val
}
