package kitty

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type KittyController struct {
	views.Controller
}

func (gc *KittyController) GetState() views.ControllerState {
	return views.LobbyControllerState
}

func (gc *KittyController) Init() {
	gc.Model = KittyModel{}
}

func (gc *KittyController) Render() string {
	playModel := gc.Model.(*KittyModel)

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

func (hc *KittyController) ParseInput(msg tea.KeyMsg) tea.Msg {
	kittyModel := hc.Model.(*KittyModel)
	switch msg.String() {
	case "right":
		kittyModel.ActiveSlotIdx++

		if kittyModel.ActiveSlotIdx > len(kittyModel.Cards)-1 {
			kittyModel.ActiveSlotIdx = 0
		}

		kittyModel.HighlighedId = kittyModel.ActiveSlotIdx //Highlighed id is to be hnalded by view
	case "left":
		kittyModel.ActiveSlotIdx--

		if kittyModel.ActiveSlotIdx < 0 {
			kittyModel.ActiveSlotIdx = len(kittyModel.Cards) - 1
		}

		kittyModel.HighlighedId = kittyModel.ActiveSlotIdx
	}

	return nil
}
func (hc *KittyController) Update(msg tea.Msg) tea.Cmd {
	return nil
}
