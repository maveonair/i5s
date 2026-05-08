package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	appincus "github.com/maveonair/i5s/internal/incus"
)

func filteredRows(rows []appincus.InstanceRow, filter string) []appincus.InstanceRow {
	filter = strings.TrimSpace(strings.ToLower(filter))
	if filter == "" {
		return rows
	}

	out := make([]appincus.InstanceRow, 0, len(rows))
	for _, row := range rows {
		if strings.Contains(row.SearchText(), filter) {
			out = append(out, row)
		}
	}
	return out
}

func renderTable(rows []appincus.InstanceRow, selected int, width int, height int, styles Styles) string {
	if height <= 0 {
		return ""
	}

	cols := columnsForWidth(width)
	var b strings.Builder
	b.WriteString(styles.TableHeader.Render(renderHeader(cols)))
	b.WriteByte('\n')

	visible := height - 1
	start := 0
	if selected >= visible && visible > 0 {
		start = selected - visible + 1
	}
	end := min(len(rows), start+visible)

	for i := start; i < end; i++ {
		line := renderRow(rows[i], cols, i == selected, (i-start)%2 == 1, styles)
		b.WriteString(line)
		b.WriteByte('\n')
	}

	if len(rows) == 0 {
		b.WriteString(styles.Muted.Render("No instances found"))
		b.WriteByte('\n')
	}

	for lipgloss.Height(b.String()) < height {
		b.WriteByte('\n')
	}

	return strings.TrimRight(b.String(), "\n")
}

func columnsForWidth(width int) []string {
	switch {
	case width >= 120:
		return []string{"NAME", "STATE", "IPV4", "IPV6", "TYPE", "SNAPSHOTS"}
	case width >= 90:
		return []string{"NAME", "STATE", "IPV4", "IPV6", "TYPE", "SNAPSHOTS"}
	case width >= 70:
		return []string{"NAME", "STATE", "IPV4", "IPV6", "TYPE"}
	default:
		return []string{"NAME", "STATE", "IPV4"}
	}
}

func renderHeader(cols []string) string {
	parts := make([]string, 0, len(cols)+1)
	parts = append(parts, " ")
	for _, col := range cols {
		parts = append(parts, pad(col, widthFor(col)))
	}
	return strings.Join(parts, " ")
}

func renderRow(row appincus.InstanceRow, cols []string, selected bool, alt bool, styles Styles) string {
	marker := " "
	if selected {
		marker = styles.SelectedMark.Render("›")
	}

	parts := []string{marker}
	for _, col := range cols {
		value := valueFor(row, col)
		text := truncate(value, widthFor(col))
		if col == "STATE" || col == "STATUS" {
			text = statusBadge(row.StatusCode, text, styles)
		}
		parts = append(parts, pad(text, widthFor(col)))
	}

	line := strings.Join(parts, " ")
	if selected {
		return styles.SelectedRow.Width(lipgloss.Width(line)).Render(line)
	}
	if alt {
		return styles.TableRowAlt.Render(line)
	}
	return styles.TableRow.Render(line)
}

func valueFor(row appincus.InstanceRow, col string) string {
	switch col {
	case "NAME":
		return row.Name
	case "TYPE":
		return row.Type
	case "STATE", "STATUS":
		return row.Status
	case "SNAPSHOTS":
		return fmt.Sprintf("%d", row.Snapshots)
	case "IPV4":
		return row.IPv4
	case "IPV6":
		return row.IPv6
	case "IMAGE":
		return row.Image
	case "PROFILES":
		return strings.Join(row.Profiles, ",")
	case "LOCATION":
		return row.Location
	case "AGE":
		return age(row.CreatedAt)
	default:
		return "-"
	}
}

func widthFor(col string) int {
	switch col {
	case "NAME":
		return 24
	case "TYPE":
		return 12
	case "STATE", "STATUS":
		return 10
	case "IPV4", "IPV6":
		return 16
	case "IMAGE":
		return 24
	case "PROFILES":
		return 18
	case "SNAPSHOTS":
		return 10
	case "LOCATION":
		return 14
	case "AGE":
		return 8
	default:
		return 10
	}
}

func pad(value string, width int) string {
	plainWidth := lipgloss.Width(value)
	if plainWidth >= width {
		return value
	}
	return value + strings.Repeat(" ", width-plainWidth)
}

func truncate(value string, width int) string {
	return truncateRunes(value, width)
}

func age(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
