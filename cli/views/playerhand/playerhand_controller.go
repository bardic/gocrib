package playerhand

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PlayerHandController struct {
	views.Controller
}

func (gc *PlayerHandController) GetState() views.ControllerState {
	return views.LoginControllerState
}

func (hc *PlayerHandController) Init() {

}

func (gc *PlayerHandController) Render() string {
	playModel := gc.Model.(*PlayerHandModel)

	s := ""
	cardViews := make([]string, 0)
	for i := 0; i < len(playModel.Cards); i++ {
		c := utils.GetCardById(playModel.Cards[i], playModel.Deck)
		view := fmt.Sprintf("%v%v", utils.GetCardSuit(c), c.Value)

		if slices.Contains(playModel.Cards, c.ID) {
			if i == playModel.HighlighedId {
				cardViews = append(cardViews, styles.SelectedFocusedStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.SelectedStyle.Render(view))
			}
		} else {
			if i == playModel.HighlighedId {
				cardViews = append(cardViews, styles.FocusedModelStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.ModelStyle.Render(view))
			}
		}
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)

	// s += styles.HelpStyle.Render(utils.BuildSubtext(v.player, v.account, utils.IsPlayerTurn(v.player.ID, v.currentTurnPlayerId)))
	return s
}

func (hc *PlayerHandController) ParseInput(msg tea.KeyMsg) tea.Msg {
	handModel := hc.Model.(*PlayerHandModel)
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		switch handModel.State {
		case queries.GamestateDiscardState:
			services.PutKitty(model.HandModifier{
				MatchId:  handModel.MatchId,
				PlayerId: handModel.PlayerId,
				CardIds:  handModel.Cards,
			})
		case queries.GamestatePlayState:
			services.PutPlay(model.HandModifier{
				MatchId:  handModel.MatchId,
				PlayerId: handModel.PlayerId,
				CardIds:  handModel.Cards,
			})
		}

		handModel.Cards = []int32{}
		return nil
	case " ":
		idx := slices.Index(handModel.HighlightedSlotIdxs, handModel.ActiveSlotIdx)

		if idx > -1 {
			handModel.HighlightedSlotIdxs = slices.Delete(handModel.HighlightedSlotIdxs, 0, 1)
		} else {
			handModel.HighlightedSlotIdxs = append(handModel.HighlightedSlotIdxs, handModel.ActiveSlotIdx)
		}
	case "right":
		handModel.ActiveSlotIdx++

		if handModel.ActiveSlotIdx > len(handModel.Cards)-1 {
			handModel.ActiveSlotIdx = 0
		}

		handModel.HighlighedId = handModel.ActiveSlotIdx //Highlighed id is to be hnalded by view
	case "left":
		handModel.ActiveSlotIdx--

		if handModel.ActiveSlotIdx < 0 {
			handModel.ActiveSlotIdx = len(handModel.Cards) - 1
		}

		handModel.HighlighedId = handModel.ActiveSlotIdx
	}

	return nil
}
func (hc *PlayerHandController) Update(msg tea.Msg) tea.Cmd {
	return nil
}
