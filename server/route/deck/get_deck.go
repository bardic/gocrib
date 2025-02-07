package deck

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"

	conn "github.com/bardic/gocrib/server/db"
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
func GetDeckByMatchId(c echo.Context) error {
	p := c.Param("matchId")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)
	ctx := c.Request().Context()

	deck, err := q.GetDeckForMatchId(ctx, &id)
	//deck, err := controller.GetDeckByMatchId(&id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cards, err := q.GetMatchCardsByTypeAndDeckId(ctx, queries.GetMatchCardsByTypeAndDeckIdParams{
		ID:    &id,
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

	return c.JSON(http.StatusOK, gamedeck)
}
