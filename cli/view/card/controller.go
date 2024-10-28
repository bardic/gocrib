package card

import (
	"cli/services"
	cliVO "cli/vo"
	"queries"
	"slices"
	"vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*cliVO.Controller
	*vo.GameMatch
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LobbyControllerState
}

func (ctrl *Controller) Init() {

}

func (ctrl *Controller) Render() string {
	cardView := ctrl.View.(*View)
	return cardView.Render()
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	cardModel := ctrl.Model.(*Model)
	cardView := ctrl.View.(*View)
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
			cardModel.SelectedCardIds,
			cardModel.CardIds[cardModel.ActiveSlotIndex])
		if idx > -1 {
			cardModel.SelectedCardIds = slices.Delete(cardModel.SelectedCardIds, 0, idx+1)
		} else {
			cardModel.SelectedCardIds = append(cardModel.SelectedCardIds, cardModel.CardIds[cardModel.ActiveSlotIndex])
		}
	case "enter":
		switch cardModel.State {
		case queries.GamestateDiscardState:
			services.PutKitty(vo.HandModifier{
				MatchId:  ctrl.ID,
				PlayerId: cardModel.LocalPlayerID,
				CardIds:  cardModel.SelectedCardIds,
			})
		case queries.GamestatePlayState:
			services.PutPlay(vo.HandModifier{
				MatchId:  ctrl.ID,
				PlayerId: cardModel.LocalPlayerID,
				CardIds:  cardModel.SelectedCardIds,
			})
		}
	}

	cardView.ActiveCardId = cardModel.ActiveSlotIndex
	cardView.SelectedCardIds = cardModel.SelectedCardIds

	return nil
}
func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (ctrl *Controller) updateActiveSlotIndex(delta int32) {
	cardModel := ctrl.Model.(*Model)
	cardModel.ActiveSlotIndex += delta

	if cardModel.ActiveSlotIndex < 0 {
		cardModel.ActiveSlotIndex = int32(len(cardModel.CardIds)) - 1
	} else if cardModel.ActiveSlotIndex > int32(len(cardModel.CardIds))-1 {
		cardModel.ActiveSlotIndex = 0
	}
}
