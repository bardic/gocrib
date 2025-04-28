package deck

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

// Returns the deck for a match id
//
//	@Summary	Get deck by match id
//	@Description
//	@Tags		deck
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"search for deck by match id"'
//	@Success	200	{object}	vo.GameDeck
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/match/{matchId}/deck/ [get]
func (h *DeckHandler) GetDeckByMatchId(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	deck, err := h.DeckStore.GetDeckForMatchId(c, &matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cards, err := h.CardStore.GetCardsForMatchIdAndState(c, queries.GetCardsForMatchIdAndStateParams{
		ID:    &matchId,
		State: "Deck",
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	gameCards := []*vo.GameCard{}
	for _, card := range cards {
		gameCards = append(gameCards, &vo.GameCard{
			Match: queries.Matchcard{
				ID:        card.Matchcardid,
				Cardid:    card.Cardid,
				Origowner: card.Origowner,
				Currowner: card.Currowner,
				State:     card.State,
			},
			Card: queries.Card{
				ID:    card.Cardid,
				Value: card.Value.Cardvalue,
				Suit:  card.Suit.Cardsuit,
				Art:   card.Art.String,
			}})
	}

	gamedeck := vo.GameDeck{
		Deck:  deck,
		Cards: gameCards,
	}

	return c.JSON(http.StatusOK, gamedeck)
}
