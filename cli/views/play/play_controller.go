package play

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PlayController struct {
	views.Controller
}

func (gc *PlayController) GetState() views.ControllerState {
	return views.LoginControllerState
}

func (gc *PlayController) Init() {
	gc.Model = PlayModel{}
}
func (gc *PlayController) Render() string {
	playModel := gc.Model.(*PlayModel)

	s := ""
	cardViews := make([]string, 0)
	for i := 0; i < len(playModel.Cards); i++ {
		c := utils.GetCardById(playModel.Cards[i], playModel.Deck)
		view := fmt.Sprintf("%v%v", utils.GetCardSuit(c), c.Value)

		if slices.Contains(v.selectedCardIds, c.ID) {
			if i == v.selectedCardId {
				cardViews = append(cardViews, styles.SelectedFocusedStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.SelectedStyle.Render(view))
			}
		} else {
			if i == v.selectedCardId {
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

func (hc *PlayController) InitView() {

}
func (hc *PlayController) ParseInput(msg tea.KeyMsg) tea.Msg {
	playModel := hc.Model.(*PlayModel)
	switch msg.String() {
	case "right":
		playModel.ActiveSlotId++

		if playModel.ActiveSlotId > len(playModel.Cards)-1 {
			playModel.ActiveSlotId = 0
		}

		playModel.HighlighedId = playModel.ActiveSlotId //Highlighed id is to be hnalded by view
	case "left":
		playModel.ActiveSlotId--

		if playModel.ActiveSlotId < 0 {
			playModel.ActiveSlotId = len(playModel.Cards) - 1
		}

		playModel.HighlighedId = playModel.ActiveSlotId
	}

	return nil
}
func (hc *PlayController) Update(msg tea.Msg) tea.Cmd {
	return nil
}
