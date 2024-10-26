package container

import (
	"strings"

	"cli/utils"
	"cli/views"

	"github.com/charmbracelet/lipgloss"
)

type ContainerView struct {
	ActiveTab int
	tabs      []views.Tab
}

func (v *ContainerView) Init() {
	v.tabs = []views.Tab{
		{
			Title: "Board",
		},
		{
			Title: "Play",
		},
		{
			Title: "Hand",
		},
		{
			Title: "Kitty",
		},
	}
}

func (v *ContainerView) Render() string {
	doc := strings.Builder{}
	renderedTabs := utils.RenderTabs(v.tabs, v.ActiveTab)
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐"))

	return doc.String()
}
