package card

import (
	"slices"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

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

func (ctrl *Controller) Render(gameMatch *vo.GameMatch) string {
	cardView := ctrl.View.(*View)
	cardView.ActiveCardId = ctrl.Model.(*Model).ActiveSlotIndex
	cardView.SelectedCardIds = ctrl.Model.(*Model).SelectedCardIds

	cardView.UIFooterVO = &vo.UIFooterVO{
		ActivePlayerId: gameMatch.Players[0].ID,
		MatchId:        gameMatch.Match.ID,
		GameState:      gameMatch.Match.Gamestate,
		LocalPlayerID:  0,
	}

	return cardView.Render(gameMatch.Players[0].Hand)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	cardModel := ctrl.Model.(*Model)
	// cardView := ctrl.View.(*View)
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
		case queries.GamestateDiscard:
			services.PutKitty(vo.HandModifier{
				// MatchId:  ctrl.ID,
				// PlayerId: cardModel.LocalPlayerID,
				CardIds: cardModel.SelectedCardIds,
			})
		case queries.GamestatePlayOwn:
			//todo pass down from CLI if player is creator
			services.PutPlay(vo.HandModifier{
				// MatchId:  ctrl.ID,
				// PlayerId: cardModel.LocalPlayerID,
				CardIds: cardModel.SelectedCardIds,
			})
		case queries.GamestatePlayOpponent:
			//todo pass down from CLI if player is opponent

			if cardModel.LocalPlayerID != ctrl.Players[0].ID {
				break
			}

			services.PutPlay(vo.HandModifier{
				// MatchId:  ctrl.ID,
				// PlayerId: cardModel.LocalPlayerID,
				CardIds: cardModel.SelectedCardIds,
			})

		}
	}

	// cardView.ActiveCardId = cardModel.ActiveSlotIndex
	// cardView.SelectedCardIds = cardModel.SelectedCardIds

	return nil
}
func (ctrl *Controller) Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd {
	ctrl.Render(gameMatch)
	return nil
}

func (ctrl *Controller) updateActiveSlotIndex(delta int) {
	cardModel := ctrl.Model.(*Model)
	cardModel.ActiveSlotIndex += delta

	if cardModel.ActiveSlotIndex < 0 {
		cardModel.ActiveSlotIndex = int(len(cardModel.CardIds)) - 1
	} else if cardModel.ActiveSlotIndex > int(len(cardModel.CardIds))-1 {
		cardModel.ActiveSlotIndex = 0
	}
}
