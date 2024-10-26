package card

import (
	"cli/views"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

type CardController struct {
	*views.Controller
}

func (cc *CardController) GetState() views.ControllerState {
	return views.LobbyControllerState
}

func (cc *CardController) Init() {

}

func (cc *CardController) Render() string {
	cardView := cc.View.(*CardView)
	return cardView.Render()
}

func (cc *CardController) ParseInput(msg tea.KeyMsg) tea.Msg {
	cardModel := cc.Model.(*CardModel)
	cardView := cc.View.(*CardView)
	switch msg.String() {
	case "right":
		cardModel.ActiveSlotIdx++

		if cardModel.ActiveSlotIdx > len(cardView.CardsToDisplay)-1 {
			cardModel.ActiveSlotIdx = 0
		}

		cardModel.HighlighedId = cardModel.ActiveSlotIdx //Highlighed id is to be hnalded by view

	case "left":
		cardModel.ActiveSlotIdx--

		if cardModel.ActiveSlotIdx < 0 {
			cardModel.ActiveSlotIdx = len(cardView.CardsToDisplay) - 1
		}

		cardModel.HighlighedId = cardModel.ActiveSlotIdx
	case " ":
		idx := slices.Index(
			cardModel.HighlightedSlotIdxs,
			cardModel.CardsToDisplay[cardModel.HighlighedId])
		if idx > -1 {
			cardModel.HighlightedSlotIdxs = slices.Delete(cardModel.HighlightedSlotIdxs, 0, 1)
		} else {
			cardModel.HighlightedSlotIdxs = append(cardModel.HighlightedSlotIdxs, cardModel.CardsToDisplay[cardModel.HighlighedId])
		}
	}

	cardView.SelectedCardId = cardModel.ActiveSlotIdx
	cardView.SelectedCardIds = cardModel.HighlightedSlotIdxs

	return nil
}
func (cc *CardController) Update(msg tea.Msg) tea.Cmd {
	return nil
}
