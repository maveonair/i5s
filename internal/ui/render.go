package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lxc/incus/v7/shared/api"
)

// KeyHelp describes one key binding shown in footer help.
type KeyHelp struct {
	Key   string
	Label string
}

func renderPanelHeight(title string, body string, width int, height int, styles Styles) string {
	if width <= 2 || height <= 2 {
		return fitToScreen(body, width, height)
	}

	innerWidth := width - 2
	contentHeight := height - 2
	titleText := " " + styles.PanelTitle.Render(title) + " "
	topFill := max(0, width-lipgloss.Width(titleText)-2)
	top := styles.Muted.Render("╭") + titleText + styles.Muted.Render(strings.Repeat("─", topFill)+"╮")
	bottom := styles.Muted.Render("╰" + strings.Repeat("─", width-2) + "╯")

	bodyLines := strings.Split(body, "\n")
	if len(bodyLines) > contentHeight {
		bodyLines = bodyLines[:contentHeight]
	}
	for len(bodyLines) < contentHeight {
		bodyLines = append(bodyLines, "")
	}

	lines := make([]string, 0, height)
	lines = append(lines, fitLine(top, width))
	for _, line := range bodyLines {
		lines = append(lines, styles.Muted.Render("│")+fitLine(line, innerWidth)+styles.Muted.Render("│"))
	}
	lines = append(lines, fitLine(bottom, width))
	return strings.Join(lines, "\n")
}

func renderKeyBar(items []KeyHelp, styles Styles) string {
	parts := make([]string, 0, len(items))
	for _, item := range items {
		parts = append(parts, renderKeyHelp(item.Key, item.Label, styles))
	}
	return strings.Join(parts, "  ")
}

func renderKeyHelp(key string, label string, styles Styles) string {
	return styles.Key.Render(key) + " " + styles.FooterText.Render(label)
}

func statusBadge(code api.StatusCode, text string, styles Styles) string {
	text = strings.ToUpper(strings.TrimSpace(text))
	switch code {
	case api.Running:
		return styles.RunningBadge.Render(text)
	case api.Stopped:
		return styles.StoppedBadge.Render(text)
	case api.Frozen:
		return styles.FrozenBadge.Render(text)
	case api.Error:
		return styles.ErrorBadge.Render(text)
	default:
		return styles.DefaultBadge.Render(text)
	}
}

func truncateRunes(value string, width int) string {
	if lipgloss.Width(value) <= width {
		return value
	}
	if width <= 1 {
		return "…"
	}
	runes := []rune(value)
	for len(runes) > 0 && lipgloss.Width(string(runes))+1 > width {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "…"
}
