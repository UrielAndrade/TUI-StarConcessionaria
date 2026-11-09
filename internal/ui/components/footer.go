package components

import (
	"fmt"
	"strings"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

type Shortcut struct {
	Key  string
	Desc string
}

func RenderFooter(shortcuts []Shortcut, statusMsg string, isConnected bool, width int) string {
	var renderedShortcuts []string
	for _, sc := range shortcuts {
		k := theme.KeyStyle.Render(fmt.Sprintf("[%s]", sc.Key))
		renderedShortcuts = append(renderedShortcuts, fmt.Sprintf("%s %s", k, sc.Desc))
	}
	left := strings.Join(renderedShortcuts, "  ")

	var connIndicator string
	if isConnected {
		connIndicator = lipgloss.NewStyle().Foreground(theme.Success).Render("● Online")
	} else {
		connIndicator = lipgloss.NewStyle().Foreground(theme.Warning).Render("○ Conectando...")
	}

	right := fmt.Sprintf("%s | %s", connIndicator, statusMsg)

	availableSpace := width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	var filler string
	if availableSpace > 0 {
		filler = strings.Repeat(" ", availableSpace)
	}

	content := lipgloss.JoinHorizontal(lipgloss.Center, left, filler, right)
	divider := lipgloss.NewStyle().Foreground(theme.BorderCol).Render(strings.Repeat("─", width))

	return lipgloss.JoinVertical(lipgloss.Left, divider, theme.FooterStyle.Width(width).Render(content))
}
