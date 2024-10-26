package container

import (
	"strings"

	"cli/utils"
	"cli/views"

	"github.com/charmbracelet/lipgloss"
)

type ContainerView struct {
	ActiveTab int
	Tabs      []views.Tab
}

func (v *ContainerView) Init() {

}

func (v *ContainerView) Render() string {
	doc := strings.Builder{}
	renderedTabs := utils.RenderTabs(v.Tabs, v.ActiveTab)
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐"))

	return doc.String()
}
