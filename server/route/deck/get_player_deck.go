package deck

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

// Returns the deck for a match playerId
//
//	@Summary	Get deck by match playerId
//	@Description
//	@Tags		deck
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"search for deck by match playerId"'
//	@Param		playerId	path		int	true	"search for deck by player playerId"'
//	@Success	200	{object}	vo.GameDeck
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/match/{matchId}/player/{playerId}/deck/ [get]
func (h *DeckHandler) GetDeckByPlayerIdAndMatchId(c echo.Context) error {
	p := c.Param("playerId")
	playerId, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	pid := c.Param("matchId")
	matchId, err := strconv.Atoi(pid)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	gameDeck, err := h.ParseQueryCardsToGameCards(c, &matchId, &playerId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, gameDeck)
}

func (h *DeckHandler) ParseQueryCardsToGameCards(ctx echo.Context, matchId, playerId *int) (*vo.GameDeck, error) {
	deck, err := h.DeckStore.GetDeckForMatchId(ctx, matchId)

	if err != nil {
		return nil, err
	}

	params := queries.GetCardsForPlayerIdFromDeckIdParams{
		Deckid:    deck.ID,
		Origowner: playerId,
	}
	cards, err := h.CardStore.GetCardsForPlayerIdFromDeckId(ctx, params)

	if err != nil {
		return nil, err
	}

	gameCards := []*vo.GameCard{}
	for _, card := range cards {
		gameCards = append(gameCards, &vo.GameCard{
			Match: queries.Matchcard{
				ID:        card.ID,
				Cardid:    card.Cardid,
				Origowner: card.Origowner,
				Currowner: card.Currowner,
				State:     card.State.Cardstate,
			},
			Card: queries.Card{
				ID:    card.Cardid,
				Value: card.Value,
				Suit:  card.Suit,
				Art:   card.Art,
			}})
	}

	gamedeck := &vo.GameDeck{
		Deck:  deck,
		Cards: gameCards,
	}

	return gamedeck, nil
}
