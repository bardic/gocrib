package views

import (
	"strings"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type GameView struct {
	GameState      model.GameState
	HighlightedIds []int
	HighlightId    int
	Cards          []model.Card
	ViewModel      ViewModel
}

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

func (v GameView) View() string {
	createInput()

	doc := strings.Builder{}

	renderedTabs := renderTabs(v.ViewModel.Tabs, v.ViewModel.ActiveTab)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	var view string
	switch v.ViewModel.GameViewState {
	case model.BoardView:
		if state.ActiveMatch != nil && state.ActiveMatch.GameState == model.CutState {
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
		view = HandView(v.HighlightId, v.HighlightedIds, v.Cards)
	case model.HandView:
		if v.GameState == 0 {
			view = "Waiting to be dealt"
		} else {
			view = HandView(v.HighlightId, v.HighlightedIds, v.Cards)
		}
	case model.KittyView:
		if len(v.Cards) == 0 {
			view = "Empty Kitty"
		} else {
			view = HandView(v.HighlightId, v.HighlightedIds, v.Cards)
		}
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(view))
	return doc.String()
}

func (v GameView) Enter() tea.Msg {
	switch v.GameState {
	case model.CutState:
		state.CutIndex = CutInput.Value()
		return services.CutDeck
	case model.DiscardState:
		p, err := utils.GetPlayerId(state.AccountId, state.ActiveMatch.Players)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}
		state.CurrentHandModifier = model.HandModifier{
			MatchId:  state.ActiveMatchId,
			PlayerId: p.Id,
			CardIds:  v.HighlightedIds,
		}
		return services.PutKitty
	}

	return nil
}

// func (s *GameView) Enter() tea.Msg {
// 	switch gameState {
// 	case model.CutState:
// 		state.CutIndex = CutInput.Value()
// 		return services.CutDeck
// 	case model.DiscardState:
// 		p, err := utils.GetPlayerId(state.AccountId, state.ActiveMatch.Players)

// 		if err != nil {
// 			utils.Logger.Sugar().Error(err)
// 		}
// 		state.CurrentHandModifier = model.HandModifier{
// 			MatchId:  state.ActiveMatchId,
// 			PlayerId: p.Id,
// 			CardIds:  highlightIds,
// 		}
// 		return services.PutKitty
// 	}

// 	return nil
// }
