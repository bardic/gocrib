package deck

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"

	conn "github.com/bardic/gocrib/server/db"
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
func GetDeckByPlayerIdAndMatchId(c echo.Context) error {
	p := c.Param("matchId")
	playerId, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// pid := c.Param("matchId")
	// playerId, err := strconv.Atoi(pid)

	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, err)
	// }

	gameDeck, err := ParseQueryCardsToGameCards(c, playerId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, gameDeck)
}

func ParseQueryCardsToGameCards(c echo.Context, playerId int) (vo.GameDeck, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)
	ctx := c.Request().Context()

	deck, err := q.GetDeckForMatchId(ctx, &playerId)
	//deck, err := controller.GetDeckByMatchId(&playerId)

	if err != nil {
		return vo.GameDeck{}, c.JSON(http.StatusInternalServerError, err)
	}

	cards, err := q.GetMatchCardsByPlayerIdAndDeckId(ctx, &playerId)

	if err != nil {
		return vo.GameDeck{}, c.JSON(http.StatusInternalServerError, err)
	}

	gameCards := []*vo.GameCard{}
	for _, card := range cards {
		gameCards = append(gameCards, &vo.GameCard{
			Match: queries.Matchcard{
				ID:        card.Matchcardid,
				Cardid:    card.Matchcardcardid,
				Origowner: card.Origowner,
				Currowner: card.Currowner,
				State:     card.State.Cardstate,
			},
			Card: queries.Card{
				ID:    card.Cardid,
				Value: card.Value.Cardvalue,
				Suit:  card.Suit.Cardsuit,
				Art:   card.Art.String,
			}})
	}

	gamedeck := vo.GameDeck{
		Deck:  &deck,
		Cards: gameCards,
	}

	return gamedeck, nil
}
