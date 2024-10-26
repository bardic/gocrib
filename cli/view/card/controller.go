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
	case "right":
		cardModel.ActiveSlotIndex++

		if cardModel.ActiveSlotIndex > len(cardView.CardsToDisplay)-1 {
			cardModel.ActiveSlotIndex = 0
		}

		cardModel.HighlighedId = cardModel.ActiveSlotIndex //Highlighed id is to be hnalded by view

	case "left":
		cardModel.ActiveSlotIndex--

		if cardModel.ActiveSlotIndex < 0 {
			cardModel.ActiveSlotIndex = len(cardView.CardsToDisplay) - 1
		}

		cardModel.HighlighedId = cardModel.ActiveSlotIndex
	case " ":
		idx := slices.Index(
			cardModel.HighlightedSlotIndexes,
			cardModel.CardsToDisplay[cardModel.HighlighedId])
		if idx > -1 {
			cardModel.HighlightedSlotIndexes = slices.Delete(cardModel.HighlightedSlotIndexes, 0, 1)
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
