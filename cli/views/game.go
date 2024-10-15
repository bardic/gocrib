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
	Cards          []model.Card
	GameTabNames   []string
	GameViewState  model.GameViewState
	ActiveSlot     int
	ActiveTab      int
	HighlighedId   int
	CutInput       textinput.Model
	gameViewInitd  bool
}

var focusedModelStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69"))

func (v GameView) Init() {
	view := &v
	if view.gameViewInitd {
		return
	}

	view.gameViewInitd = true

	view.GameTabNames = []string{"Board", "Play", "Hand", "Kitty"}

	view.CutInput = textinput.New()
	view.CutInput.Placeholder = "0"
	view.CutInput.CharLimit = 5
	view.CutInput.Width = 5
}

func (v GameView) View() string {
	doc := strings.Builder{}

	renderedTabs := renderTabs(v.GameTabNames, v.ActiveTab)

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row)
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, "───────────────────────────────────────────────────────────────────┐")
	doc.WriteString(row)
	doc.WriteString("\n")

	var view string
	switch v.GameViewState {
	case model.BoardView:
		if state.ActiveMatch != nil && state.ActiveMatch.GameState == model.CutState {
			v.CutInput.Focus()
			view = v.CutInput.View() + " \n"
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
		view = HandView(v.HighlighedId, v.HighlightedIds, v.Cards)
	case model.HandView:
		if v.GameState == 0 {
			view = "Waiting to be dealt"
		} else {
			view = HandView(v.HighlighedId, v.HighlightedIds, v.Cards)
		}
	case model.KittyView:
		if len(v.Cards) == 0 {
			view = "Empty Kitty"
		} else {
			view = HandView(v.HighlighedId, v.HighlightedIds, v.Cards)
		}
	}

	doc.WriteString(styles.WindowStyle.Width(100).Render(view))
	return doc.String()
}

func (v GameView) Enter() tea.Msg {
	switch v.GameState {
	case model.CutState:
		state.CutIndex = v.CutInput.Value()
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

func (v GameView) Update(msg tea.Msg) tea.Cmd {
	v.Init()
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
