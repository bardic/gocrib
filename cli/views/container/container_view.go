package container

import (
	"strings"

	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/charmbracelet/lipgloss"
)

type ContainerView struct {
	activeTab int
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
	renderedTabs := utils.RenderTabs(v.tabs, v.activeTab)
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐"))

	return doc.String()
}
