package components

import (
	"fmt"
	"strings"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

type TabItem struct {
	Key   string
	Title string
}

func RenderHeader(activeTab int, tabs []TabItem, width int) string {
	logo := theme.HeaderTitleStyle.Render("★ STAR CONCESSIONÁRIA")
	version := theme.HeaderSubStyle.Render("v1.0.0")

	var renderedTabs []string
	for i, t := range tabs {
		label := fmt.Sprintf("[%s] %s", t.Key, t.Title)
		if i == activeTab {
			renderedTabs = append(renderedTabs, theme.TabActiveStyle.Render(label))
		} else {
			renderedTabs = append(renderedTabs, theme.TabInactiveStyle.Render(label))
		}
	}
	tabsRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	topBar := lipgloss.JoinHorizontal(lipgloss.Center, logo, " ", version)
	availableSpace := width - lipgloss.Width(topBar) - lipgloss.Width(tabsRow) - 2
	var filler string
	if availableSpace > 0 {
		filler = strings.Repeat(" ", availableSpace)
	}

	header := lipgloss.JoinHorizontal(lipgloss.Center, topBar, filler, tabsRow)
	divider := lipgloss.NewStyle().Foreground(theme.BorderCol).Render(strings.Repeat("─", width))

	return lipgloss.JoinVertical(lipgloss.Left, header, divider)
}
