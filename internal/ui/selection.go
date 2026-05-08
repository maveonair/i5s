package ui

import appincus "github.com/maveonair/i5s/internal/incus"

func (m Model) selected() appincus.InstanceRow {
	rows := filteredRows(m.instances, m.filter)
	if len(rows) == 0 || m.selectedIndex < 0 || m.selectedIndex >= len(rows) {
		return appincus.InstanceRow{}
	}
	return rows[m.selectedIndex]
}

func (m *Model) moveSelection(delta int) {
	rows := filteredRows(m.instances, m.filter)
	if len(rows) == 0 {
		m.selectedIndex = 0
		m.selectedName = ""
		return
	}
	m.selectedIndex += delta
	if m.selectedIndex < 0 {
		m.selectedIndex = 0
	}
	if m.selectedIndex >= len(rows) {
		m.selectedIndex = len(rows) - 1
	}
	m.selectedName = rows[m.selectedIndex].Name
}

func (m *Model) restoreSelection() {
	rows := filteredRows(m.instances, m.filter)
	if len(rows) == 0 {
		m.selectedIndex = 0
		m.selectedName = ""
		return
	}
	for i, row := range rows {
		if row.Name == m.selectedName {
			m.selectedIndex = i
			return
		}
	}
	if m.selectedIndex >= len(rows) {
		m.selectedIndex = len(rows) - 1
	}
	m.selectedName = rows[m.selectedIndex].Name
}
