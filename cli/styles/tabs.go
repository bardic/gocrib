package styles

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/charmbracelet/lipgloss"
)

func RenderTabs(tabs []cliVO.Tab, activeTab int) []string {
	var renderedTabs []string

	for i, t := range tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(tabs)-1, i == activeTab
		if isActive {
			style = ActiveTabStyle
		} else {
			style = InactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "└"
		} else if isLast && !isActive {
			border.BottomRight = "┴"
		}

		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t.Title))
	}

	return renderedTabs
}
