package card

import (
	"slices"

	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	model *Model
	view  *View
}

func NewController(name string, match *vo.Match, player *vo.Player) *Controller {
	ctrl := &Controller{
		model: &Model{
			ActiveSlotIndex: 0,
			SelectedCardIDs: []int{},
			State:           match.Gamestate,
			HandVO:          &cliVO.HandVO{},
			LocalPlayer:     player,
			ActivePlayerID:  match.Dealerid,
			Name:            name,
			GameMatchID:     match.ID,
		},
	}

	v := NewCardView(player, name)

	ctrl.view = v

	return ctrl
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LobbyControllerState
}

func (ctrl *Controller) GetName() string {
	return ctrl.model.Name
}

func (ctrl *Controller) Render(gameMatch *vo.Match, gameDeck *vo.Deck) string {
	ctrl.model.State = gameMatch.Gamestate
	ctrl.model.Deck = gameDeck
	ctrl.model.ActivePlayerID = gameMatch.Dealerid

	// p := utils.GetPlayerForAccountID(ctrl.model.LocalPlayer.Accountid, gameMatch)

	var hand []int
	// switch ctrl.model.Name {
	// case "Play":
	// 	hand = utils.IDFromCards(p.Play)
	// case "Hand":
	// 	hand = utils.IDFromCards(p.Hand)
	// case "Kitty":
	// 	hand = utils.IDFromCards(p.Kitty)
	// default:
	// 	return "Error: Unknown tabname"
	// }

	ctrl.model.CardIDs = hand

	return ctrl.view.Render(gameMatch, gameDeck, hand)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	// Highlight card to the right
	case "right":
		ctrl.updateActiveSlotIndex(1)
	// Highlight card to the left
	case "left":
		ctrl.updateActiveSlotIndex(-1)
	// Select card
	case " ":
		idx := slices.Index(
			ctrl.model.SelectedCardIDs,
			ctrl.model.CardIDs[ctrl.model.ActiveSlotIndex])
		if idx > -1 {
			ctrl.model.SelectedCardIDs = slices.Delete(ctrl.model.SelectedCardIDs, 0, idx+1)
		} else {
			ctrl.model.SelectedCardIDs = append(ctrl.model.SelectedCardIDs, ctrl.model.CardIDs[ctrl.model.ActiveSlotIndex])
		}

		ctrl.view.SelectedCardIDs = ctrl.model.SelectedCardIDs
	case "enter":
		switch ctrl.model.State {
		// case queries.GamestateDiscard:
		// 	services.PutKitty(
		// 		ctrl.model.GameMatchID,
		// 		ctrl.model.LocalPlayer.ID,
		// 		ctrl.model.ActivePlayerID,
		// 		vo.HandModifier{
		// 			CardIDs: ctrl.model.SelectedCardIDs,
		// 		},
		// 	)

		// 	ctrl.model.SelectedCardIDs = []int{}
		// case queries.GamestatePlay:
		// 	services.PutPlay(
		// 		ctrl.model.GameMatchID,
		// 		ctrl.model.LocalPlayer.ID,
		// 		ctrl.model.ActivePlayerID,
		// 		vo.HandModifier{
		// 			CardIDs: ctrl.model.SelectedCardIDs,
		// 		},
		// 	)
		//
		// 	ctrl.model.SelectedCardIDs = []int{}
		}
	}

	return nil
}

func (ctrl *Controller) Update(_ tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	// ctrl.view.Update(msg)
	return cmd
}

func (ctrl *Controller) updateActiveSlotIndex(delta int) {
	cardModel := ctrl.model
	cardModel.ActiveSlotIndex += delta

	if cardModel.ActiveSlotIndex < 0 {
		cardModel.ActiveSlotIndex = len(cardModel.CardIDs) - 1
	} else if cardModel.ActiveSlotIndex > len(cardModel.CardIDs)-1 {
		cardModel.ActiveSlotIndex = 0
	}

	ctrl.view.ActiveCardID = cardModel.ActiveSlotIndex
}
