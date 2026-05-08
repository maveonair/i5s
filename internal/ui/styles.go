package ui

import "github.com/charmbracelet/lipgloss"

// Styles groups the Lipgloss styles used by the TUI renderer.
type Styles struct {
	AppTitle     lipgloss.Style
	TopBar       lipgloss.Style
	ContextLabel lipgloss.Style
	ContextValue lipgloss.Style
	PanelTitle   lipgloss.Style
	TableHeader  lipgloss.Style
	TableRow     lipgloss.Style
	TableRowAlt  lipgloss.Style
	SelectedRow  lipgloss.Style
	SelectedMark lipgloss.Style
	Key          lipgloss.Style
	FooterText   lipgloss.Style
	Modal        lipgloss.Style
	ModalDanger  lipgloss.Style
	Error        lipgloss.Style
	Warning      lipgloss.Style
	Success      lipgloss.Style
	Muted        lipgloss.Style
	RunningBadge lipgloss.Style
	StoppedBadge lipgloss.Style
	FrozenBadge  lipgloss.Style
	ErrorBadge   lipgloss.Style
	DefaultBadge lipgloss.Style
}

// NewStyles returns the default visual theme for i5s.
func NewStyles() Styles {
	border := lipgloss.Color("240")
	accent := lipgloss.Color("80")
	muted := lipgloss.Color("242")
	panelBg := lipgloss.Color("235")
	selectedBg := lipgloss.Color("24")
	danger := lipgloss.Color("203")

	return Styles{
		AppTitle:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("231")).Background(accent).Padding(0, 1),
		TopBar:       lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(border).Padding(0, 1),
		ContextLabel: lipgloss.NewStyle().Foreground(muted),
		ContextValue: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("87")),
		PanelTitle:   lipgloss.NewStyle().Bold(true).Foreground(accent),
		TableHeader:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("250")),
		TableRow:     lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		TableRowAlt:  lipgloss.NewStyle().Foreground(lipgloss.Color("247")),
		SelectedRow:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("231")).Background(selectedBg),
		SelectedMark: lipgloss.NewStyle().Bold(true).Foreground(accent),
		Key:          lipgloss.NewStyle().Foreground(lipgloss.Color("231")).Background(lipgloss.Color("238")).Padding(0, 1),
		FooterText:   lipgloss.NewStyle().Foreground(muted),
		Modal:        lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(border).Background(panelBg).Padding(1, 2),
		ModalDanger:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(danger).Background(panelBg).Padding(1, 2),
		Error:        lipgloss.NewStyle().Bold(true).Foreground(danger),
		Warning:      lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		Success:      lipgloss.NewStyle().Foreground(lipgloss.Color("114")),
		Muted:        lipgloss.NewStyle().Foreground(muted),
		RunningBadge: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("114")).Background(lipgloss.Color("22")).Padding(0, 1),
		StoppedBadge: lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Background(lipgloss.Color("238")).Padding(0, 1),
		FrozenBadge:  lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("229")).Background(lipgloss.Color("94")).Padding(0, 1),
		ErrorBadge:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("231")).Background(danger).Padding(0, 1),
		DefaultBadge: lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Background(lipgloss.Color("238")).Padding(0, 1),
	}
}
