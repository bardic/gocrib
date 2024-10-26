package card

import (
	cliVO "cli/vo"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*cliVO.Controller
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
			cardModel.HighlightedSlotIndexes,
			cardModel.CardsToDisplay[cardModel.HighlighedId])
		if idx > -1 {
			cardModel.HighlightedSlotIndexes = slices.Delete(cardModel.HighlightedSlotIndexes, idx, idx+1)
		} else {
			cardModel.HighlightedSlotIndexes = append(cardModel.HighlightedSlotIndexes, cardModel.CardsToDisplay[cardModel.HighlighedId])
		}
	}

	cardView.SelectedCardId = cardModel.ActiveSlotIndex
	cardView.SelectedCardIds = cardModel.HighlightedSlotIndexes

	return nil
}
func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	return nil
}

func (ctrl *Controller) updateActiveSlotIndex(delta int) {
	cardModel := ctrl.Model.(*Model)

	cardModel.ActiveSlotIndex += delta

	if cardModel.ActiveSlotIndex < 0 {
		cardModel.ActiveSlotIndex = len(cardModel.CardsToDisplay) - 1
	} else if cardModel.ActiveSlotIndex > len(cardModel.CardsToDisplay)-1 {
		cardModel.ActiveSlotIndex = 0
	}

	cardModel.HighlighedId = cardModel.ActiveSlotIndex
}
