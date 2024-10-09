package views

import (
	"strings"

	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var focusedModelStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69"))

var CutInput textinput.Model
var initd bool

func createInput() {
	if initd {
		return
	}

	initd = true

	CutInput = textinput.New()
	CutInput.Placeholder = "0"
	CutInput.CharLimit = 5
	CutInput.Width = 5

}

func GameView(highlightId int, highlightedIds []int, cards []model.Card, m ViewModel, s model.GameState) string {

	createInput()

	doc := strings.Builder{}

	renderedTabs := renderTabs(m.Tabs, m.ActiveTab)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	var view string
	switch m.GameViewState {
	case model.BoardView:
		if state.ActiveMatch.GameState == model.CutState {
			CutInput.Focus()
			view = CutInput.View() + " \n"
		} else {
			view = "\n"
		}

		view = view + `
	○•○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	○•○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	----- ----- ----- ----- ----- ----- ----- ----- ----- -----		
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○		
	----- ----- ----- ----- ----- ----- ----- ----- ----- -----		
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○
	○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○○○○○ ○
`

		s := lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(view), focusedModelStyle.Render("59"))
		view = s
	case model.PlayView:
		view = HandView(highlightId, highlightedIds, cards)
	case model.HandView:
		if s == 0 {
			view = "Waiting to be dealt"
		} else {
			view = HandView(highlightId, highlightedIds, cards)
		}
	case model.KittyView:
		if len(cards) == 0 {
			view = "Empty Kitty"
		} else {
			view = HandView(highlightId, highlightedIds, cards)
		}
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(view))
	return doc.String()
}
