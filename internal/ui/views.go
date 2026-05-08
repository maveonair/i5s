package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current TUI screen.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	base := fitToScreen(m.instancesView(), m.width, m.height)
	switch m.view {
	case ViewLogs, ViewConsole, ViewHelp:
		return fitToScreen(m.fullscreenView(), m.width, m.height)
	case ViewRemotePicker, ViewProjectPicker, ViewLogPicker:
		return fitToScreen(overlay(base, m.pickerView(), m.width, m.height), m.width, m.height)
	case ViewConfirmStop, ViewConfirmDelete, ViewConfirmClearConsole:
		return fitToScreen(overlay(base, m.confirmView(), m.width, m.height), m.width, m.height)
	default:
		return fitToScreen(base, m.width, m.height)
	}
}

func (m Model) instancesView() string {
	rows := filteredRows(m.instances, m.filter)
	header := m.headerView(len(rows))
	footer := m.footerView()

	// The main screen is intentionally header -> body -> footer. The table panel
	// owns all remaining vertical space so the resource list feels like a real
	// fullscreen TUI rather than a compact block of text.
	bodyHeight := max(3, m.height-lipgloss.Height(header)-lipgloss.Height(footer))
	tableHeight := max(1, bodyHeight-2)
	title := fmt.Sprintf("Instances %d", len(rows))
	if m.filter != "" {
		title = fmt.Sprintf("Instances %d/%d", len(rows), len(m.instances))
	}
	table := renderPanelHeight(title, renderTable(rows, m.selectedIndex, max(20, m.width-4), tableHeight, m.styles), m.width, bodyHeight, m.styles)

	parts := []string{header, table}
	parts = append(parts, footer)
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) headerView(count int) string {
	refreshed := "never"
	if !m.lastRefresh.IsZero() {
		refreshed = m.lastRefresh.Format("15:04:05")
	}

	filter := "all"
	if m.filter != "" || m.filtering {
		filter = m.filter
		if filter == "" {
			filter = "typing"
		}
	}

	parts := []string{
		m.styles.AppTitle.Render("i5s"),
		m.styles.PanelTitle.Render("INSTANCES"),
		m.styles.ContextLabel.Render("remote") + " " + m.styles.ContextValue.Render(m.service.CurrentRemote()),
		m.styles.ContextLabel.Render("project") + " " + m.styles.ContextValue.Render(m.service.CurrentProject()),
		m.styles.ContextLabel.Render("items") + " " + m.styles.ContextValue.Render(fmt.Sprintf("%d", count)),
		m.styles.ContextLabel.Render("filter") + " " + m.styles.ContextValue.Render(filter),
		m.styles.ContextLabel.Render("refreshed") + " " + m.styles.ContextValue.Render(refreshed),
	}
	if m.operation != "" {
		parts = append(parts, m.styles.Warning.Render(m.operation))
	}
	return m.styles.TopBar.Width(max(1, m.width-4)).Render(strings.Join(parts, "  "))
}

func (m Model) footerView() string {
	if m.errorMessage != "" {
		return m.styles.Error.Render("error ") + m.styles.FooterText.Render(m.errorMessage)
	}
	if m.statusMessage != "" {
		return m.styles.Success.Render("status ") + m.styles.FooterText.Render(m.statusMessage)
	}
	if m.filtering {
		return renderKeyHelp("filter", m.filter, m.styles)
	}
	return renderKeyBar([]KeyHelp{{"enter", "shell"}, {"e", "config"}, {"l", "logs"}, {"c", "console"}, {"s", "stop"}, {"S", "start"}, {"d", "delete"}, {"R", "remote"}, {"p", "project"}, {"?", "help"}}, m.styles)
}

func (m Model) fullscreenView() string {
	title := "Help"
	meta := ""
	footerItems := []KeyHelp{{"esc", "back"}, {"q", "back"}}
	if m.view == ViewLogs {
		title = "LOGS"
		meta = m.styles.ContextLabel.Render("instance") + " " + m.styles.ContextValue.Render(m.selected().Name) + "  " + m.styles.ContextLabel.Render("file") + " " + m.styles.ContextValue.Render(m.selectedLog)
		footerItems = []KeyHelp{{"j/k", "scroll"}, {"g", "top"}, {"G", "bottom"}, {"r", "refresh"}, {"esc", "back"}}
	} else if m.view == ViewConsole {
		title = "CONSOLE"
		meta = m.styles.ContextLabel.Render("instance") + " " + m.styles.ContextValue.Render(m.selected().Name)
		footerItems = []KeyHelp{{"j/k", "scroll"}, {"r", "refresh"}, {"d", "clear"}, {"esc", "back"}}
	} else {
		title = "HELP"
	}

	headerContent := strings.Join([]string{m.styles.AppTitle.Render("i5s"), m.styles.PanelTitle.Render(title), meta}, "  ")
	header := m.styles.TopBar.Width(max(1, m.width-4)).Render(headerContent)
	footer := renderKeyBar(footerItems, m.styles)
	m.viewport.Width = m.width
	m.viewport.Height = max(1, m.height-lipgloss.Height(header)-lipgloss.Height(footer))
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
		footer,
	)
}

func (m Model) pickerView() string {
	var title string
	var lines []string
	switch m.view {
	case ViewRemotePicker:
		title = "Switch Remote"
		lines = append(lines, m.styles.TableHeader.Render(fmt.Sprintf("%-18s %-32s %-10s %s", "REMOTE", "URL", "PROTOCOL", "PROJECT")))
		for i, row := range m.remotes {
			line := fmt.Sprintf("%-18s %-32s %-10s %s", truncateRunes(row.Name, 18), truncateRunes(row.Addr, 32), row.Protocol, row.DefaultProject)
			lines = append(lines, pickerLine(line, i == m.pickerIndex, m.styles))
		}
	case ViewProjectPicker:
		title = "Switch Project"
		lines = append(lines, m.styles.TableHeader.Render(fmt.Sprintf("%-18s %-36s %s", "PROJECT", "DESCRIPTION", "USED BY")))
		for i, row := range m.projects {
			line := fmt.Sprintf("%-18s %-36s %d", truncateRunes(row.Name, 18), truncateRunes(row.Description, 36), row.UsedBy)
			lines = append(lines, pickerLine(line, i == m.pickerIndex, m.styles))
		}
	case ViewLogPicker:
		title = "Select Log"
		for i, logName := range m.logFiles {
			lines = append(lines, pickerLine(logName, i == m.pickerIndex, m.styles))
		}
	}
	if len(lines) == 0 {
		lines = append(lines, m.styles.Muted.Render("No entries"))
	}
	content := m.styles.PanelTitle.Render(title) + "\n\n" + strings.Join(lines, "\n") + "\n\n" + renderKeyBar([]KeyHelp{{"enter", "select"}, {"esc", "cancel"}}, m.styles)
	return m.styles.Modal.Render(content)
}

func (m Model) confirmView() string {
	selected := m.selected()
	var text string
	switch m.view {
	case ViewConfirmStop:
		text = m.styles.PanelTitle.Render("Stop Instance") + "\n\n" + fmt.Sprintf("Stop instance %s?", m.styles.ContextValue.Render(selected.Name))
	case ViewConfirmDelete:
		text = m.styles.Error.Render("Delete Instance") + "\n\n" + fmt.Sprintf("Delete instance %s?\n%s", m.styles.ContextValue.Render(selected.Name), m.styles.Error.Render("This cannot be undone."))
	case ViewConfirmClearConsole:
		text = m.styles.PanelTitle.Render("Clear Console Log") + "\n\n" + fmt.Sprintf("Clear console log for %s?", m.styles.ContextValue.Render(selected.Name))
	}
	content := text + "\n\n" + renderKeyBar([]KeyHelp{{"y", "confirm"}, {"n", "cancel"}}, m.styles)
	if m.view == ViewConfirmDelete {
		return m.styles.ModalDanger.Render(content)
	}
	return m.styles.Modal.Render(content)
}

func pickerLine(line string, selected bool, styles Styles) string {
	if selected {
		return styles.SelectedRow.Render("› " + line)
	}
	return "  " + line
}

func overlay(base, modal string, width, height int) string {
	baseLines := strings.Split(fitToScreen(base, width, height), "\n")
	modalLines := strings.Split(modal, "\n")
	modalWidth := 0
	for _, line := range modalLines {
		modalWidth = max(modalWidth, lipgloss.Width(line))
	}

	top := max(0, (height-len(modalLines))/2)
	left := max(0, (width-modalWidth)/2)

	for i, line := range modalLines {
		idx := top + i
		if idx >= len(baseLines) {
			break
		}
		prefix := strings.Repeat(" ", left)
		baseLines[idx] = prefix + line
	}
	return strings.Join(baseLines, "\n")
}

func fitToScreen(view string, width, height int) string {
	if width <= 0 || height <= 0 {
		return view
	}

	lines := strings.Split(view, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}

	for len(lines) < height {
		lines = append(lines, "")
	}

	for i, line := range lines {
		lines[i] = fitLine(line, width)
	}

	return strings.Join(lines, "\n")
}

func fitLine(line string, width int) string {
	line = lipgloss.NewStyle().MaxWidth(width).Render(line)
	lineWidth := lipgloss.Width(line)
	if lineWidth >= width {
		return line
	}
	return line + strings.Repeat(" ", width-lineWidth)
}
