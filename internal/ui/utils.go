package ui

func (m model) getCurrentServer() string {
	return m.table.SelectedRow()[0]
}
