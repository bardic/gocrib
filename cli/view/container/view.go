package container

import (
	"strings"

	"github.com/bardic/gocrib/cli/styles"
	cliVO "github.com/bardic/gocrib/cli/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	ActiveTab int
	Tabs      map[int]cliVO.IGameController
}

func NewView(activeTab int, tabs map[int]cliVO.IGameController) *View {
	return &View{
		ActiveTab: activeTab,
		Tabs:      tabs,
	}
}

func (view *View) Init() {
}

func (view *View) Render() string {
	doc := strings.Builder{}
	renderedTabs := styles.RenderTabs([]string{"Board", "Play", "Hand", "Kitty"}, view.ActiveTab)
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
