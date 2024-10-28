package container

import (
	"strings"

	"cli/utils"
	cliVO "cli/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	ActiveTab int
	Tabs      []cliVO.Tab
}

func (view *View) Init() {

}

func (view *View) Render() string {
	doc := strings.Builder{}
	renderedTabs := utils.RenderTabs(view.Tabs, view.ActiveTab)
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Bottom, "──────────────────────────────────────────┐"))

	return doc.String()
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	return ""
}
