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
	*cliVO.Controller
	*vo.GameMatch
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LobbyControllerState
}

func (ctrl *Controller) Init() {

}

func (ctrl *Controller) Render(gameMatch *vo.GameMatch) string {
	model := ctrl.Model.(*Model)
	cardView := ctrl.View.(*View)
	cardView.ActiveCardId = model.ActiveSlotIndex
	cardView.SelectedCardIds = model.SelectedCardIds
	cardView.Deck = model.Deck

	localPlayerId := model.LocalPlayerID

	cardView.UIFooterVO = &vo.UIFooterVO{
		ActivePlayerId: gameMatch.Currentplayerturn,
		MatchId:        gameMatch.Match.ID,
		GameState:      gameMatch.Match.Gamestate,
		LocalPlayerID:  &localPlayerId,
	}

	var localPlayer vo.GamePlayer
	for _, player := range ctrl.Players {
		if *player.ID == model.LocalPlayerID {
			localPlayer = player
		}
	}

	ids := utils.IdFromCards(localPlayer.Hand)

	ctrl.Model.(*Model).CardIds = ids

	return cardView.Render(ids)
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
			services.PutKitty(
				ctrl.ID,
				&cardModel.LocalPlayerID,
				vo.HandModifier{
					CardIds: cardModel.SelectedCardIds,
				},
			)

			//fmt.Println(msg)
		case queries.GamestatePlay:
			services.PutPlay(
				ctrl.ID,
				&cardModel.LocalPlayerID,
				vo.HandModifier{
					CardIds: cardModel.SelectedCardIds,
				},
			)
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
