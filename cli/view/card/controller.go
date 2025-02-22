package card

import (
	"slices"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	model *Model
	view  *View
}

func NewController(name string, match *vo.GameMatch, player *vo.GamePlayer) *Controller {

	ctrl := &Controller{
		model: &Model{
			ActiveSlotIndex: 0,
			SelectedCardIds: []int{},
			State:           match.Match.Gamestate,
			HandVO:          &cliVO.HandVO{},
			LocalPlayer:     player,
			Name:            name,
			GameMatchId:     match.Match.ID,
		},
	}

	v := NewCardView(match, player, ctrl.model.HandVO.Deck, name)

	ctrl.view = v

	return ctrl
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LobbyControllerState
}

func (ctrl *Controller) GetName() string {
	return ctrl.model.Name
}

func (ctrl *Controller) Render(gameMatch *vo.GameMatch, gameDeck *vo.GameDeck) string {
	ctrl.model.State = gameMatch.Match.Gamestate
	ctrl.model.HandVO.Deck = gameDeck

	p := utils.GetPlayerForAccountId(ctrl.model.LocalPlayer.Accountid, gameMatch)

	var hand []int
	switch ctrl.model.Name {
	case "Play":
		hand = utils.IdFromCards(p.Play)
	case "Hand":
		hand = utils.IdFromCards(p.Hand)
	case "Kitty":
		hand = utils.IdFromCards(p.Kitty)
	default:
		return "Error: Unknown tabname"
	}

	ctrl.model.CardIds = hand

	return ctrl.view.Render(gameMatch, gameDeck, hand)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {

	switch msg.String() {
	//Highlight card to the right
	case "right":
		ctrl.updateActiveSlotIndex(1)
	//Highlight card to the left
	case "left":
		ctrl.updateActiveSlotIndex(-1)
	//Select card
	case " ":
		idx := slices.Index(
			ctrl.model.SelectedCardIds,
			ctrl.model.CardIds[ctrl.model.ActiveSlotIndex])
		if idx > -1 {
			ctrl.model.SelectedCardIds = slices.Delete(ctrl.model.SelectedCardIds, 0, idx+1)
		} else {
			ctrl.model.SelectedCardIds = append(ctrl.model.SelectedCardIds, ctrl.model.CardIds[ctrl.model.ActiveSlotIndex])
		}

		ctrl.view.SelectedCardIds = ctrl.model.SelectedCardIds
	case "enter":
		switch ctrl.model.State {
		case queries.GamestateDiscard:
			services.PutKitty(
				ctrl.model.GameMatchId,
				ctrl.model.LocalPlayer.ID,
				vo.HandModifier{
					CardIds: ctrl.model.SelectedCardIds,
				},
			)

			ctrl.model.SelectedCardIds = []int{}
		case queries.GamestatePlay:
			services.PutPlay(
				ctrl.model.GameMatchId,
				ctrl.model.LocalPlayer.ID,
				vo.HandModifier{
					CardIds: ctrl.model.SelectedCardIds,
				},
			)

			ctrl.model.SelectedCardIds = []int{}
		}
	}

	return nil
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	//ctrl.view.Update(msg)
	return cmd
}

func (ctrl *Controller) updateActiveSlotIndex(delta int) {
	cardModel := ctrl.model
	cardModel.ActiveSlotIndex += delta

	if cardModel.ActiveSlotIndex < 0 {
		cardModel.ActiveSlotIndex = int(len(cardModel.CardIds)) - 1
	} else if cardModel.ActiveSlotIndex > int(len(cardModel.CardIds))-1 {
		cardModel.ActiveSlotIndex = 0
	}

	ctrl.view.ActiveCardId = cardModel.ActiveSlotIndex
}
