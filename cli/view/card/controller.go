package card

import (
	"encoding/json"
	"slices"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/game"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*game.Controller
	*vo.GameMatch
	tabName string
}

func NewController(name string, match *vo.GameMatch, player *vo.GamePlayer) *Controller {

	ctrl := &Controller{
		Controller: &game.Controller{},
		GameMatch:  match,
		tabName:    name,
	}

	handModel := ctrl.getHandModelForCardIds(
		*player.ID,
		*match.ID,
		utils.IdFromCards(player.Play),
	)

	m := &Model{
		ViewModel: &cliVO.ViewModel{
			Name: name,
		},
		ActiveSlotIndex: 0,
		SelectedCardIds: []int{},
		HandVO:          handModel,
		State:           match.Match.Gamestate,
	}

	v := NewCardView(match, player, handModel.Deck)

	ctrl.Model = m
	ctrl.View = v

	return ctrl
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LobbyControllerState
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

	var localPlayer *vo.GamePlayer
	for _, player := range ctrl.Players {
		if *player.ID == model.LocalPlayerID {
			localPlayer = player
		}
	}

	var cardIds []int
	switch ctrl.tabName {
	case "Play":
		cardIds = utils.IdFromCards(localPlayer.Play)
	case "Hand":
		cardIds = utils.IdFromCards(localPlayer.Hand)
	case "Kitty":
		cardIds = utils.IdFromCards(localPlayer.Kitty)
	}

	ctrl.Model.(*Model).CardIds = cardIds

	return cardView.Render(cardIds)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	cardModel := ctrl.Model.(*Model)

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
			activePlayer := utils.GetPlayerForAccountId(&cardModel.LocalPlayerID, ctrl.GameMatch)
			if activePlayer.Isready {
				return nil
			}
			services.PutKitty(
				ctrl.ID,
				&cardModel.LocalPlayerID,
				vo.HandModifier{
					CardIds: cardModel.SelectedCardIds,
				},
			)
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

func (ctrl *Controller) getHandModelForCardIds(localPlayerId, matchId int, cardIds []int) *cliVO.HandVO {
	gameDeck := ctrl.getDeckByPlayerIdAndMatchId(localPlayerId, matchId)

	handModel := &cliVO.HandVO{
		LocalPlayerID: localPlayerId,
		CardIds:       cardIds,
		Deck:          gameDeck,
	}

	return handModel
}

func (ctrl *Controller) getDeckByPlayerIdAndMatchId(playerId, matchId int) *vo.GameDeck {
	var deck *vo.GameDeck

	resp := services.GetDeckByPlayIdAndMatchId(playerId, matchId)
	err := json.Unmarshal(resp.([]byte), &deck)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	return deck
}
