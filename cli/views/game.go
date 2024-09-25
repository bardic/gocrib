package views

import (
	"strings"

	"github.com/bardic/cribbagev2/cli/styles"
	"github.com/bardic/cribbagev2/model"
	"github.com/charmbracelet/lipgloss"
)

func GameView(highlightId int, highlightedIds []int, cards []model.Card, m ViewModel) string {
	doc := strings.Builder{}
	var renderedTabs []string

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.ActiveTab
		if isActive {
			style = styles.ActiveTabStyle
		} else {
			style = styles.InactiveTabStyle
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
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	var view string
	switch m.ViewState {
	case model.BoardView:
		view = ""
	case model.PlayView:
		view = HandView(highlightId, highlightedIds, cards)
	case model.HandView:
		view = HandView(highlightId, highlightedIds, cards)
	case model.KittyView:
		view = HandView(highlightId, highlightedIds, cards)
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(view))
	return doc.String()
}
