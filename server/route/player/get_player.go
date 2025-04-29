package player

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	Get player by playerId
//	@Description
//	@Tags		players
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"search for player in match"
//	@Param		playerId	path		int	true	"for id"
//	@Success	200	{object}	vo.GamePlayer
//	@Failure	400	{object}	error
//	@Failure	404	{object}	error
//	@Failure	500	{object}	error
//	@Router		/match/{matchId}/player/{playerId} [get]

func (h *PlayerHandler) GetPlayer(c echo.Context) error {
	playerId, err := strconv.Atoi(c.Param("playerId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	matchId, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	player, err := h.PlayerStore.GetPlayerById(c, &playerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cardStates := []queries.Cardstate{
		queries.CardstateDeck,
		queries.CardstateHand,
		queries.CardstateKitty,
	}

	gameHands := make(map[queries.Cardstate][]vo.GameCard, 3)

	for _, cardState := range cardStates {

		hand, err := h.CardStore.GetCardsForMatchIdAndState(c, queries.GetCardsForMatchIdAndStateParams{
			ID:    &matchId,
			State: cardState,
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		cards := make([]vo.GameCard, len(hand))
		for i, res := range hand {
			cards[i] = vo.GameCard{
				Match: queries.Matchcard{
					ID:        res.Matchcardid,
					Cardid:    res.Cardid,
					Origowner: res.Origowner,
					Currowner: res.Currowner,
				},
				Card: queries.Card{
					ID:    res.Cardid,
					Value: res.Value.Cardvalue,
					Suit:  res.Suit.Cardsuit,
					Art:   res.Art.String,
				},
			}
		}

		gameHands[cardState] = cards
	}

	gamePlayer := vo.GamePlayer{
		Player: player,
		Hand:   gameHands[queries.CardstateHand],
		Play:   gameHands[queries.CardstatePlay],
		Kitty:  gameHands[queries.CardstateKitty],
	}

	return c.JSON(http.StatusOK, gamePlayer)
}
