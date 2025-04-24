package deck

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	conn "github.com/bardic/gocrib/server/db"

	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	PutSHuffle
//	@Description
//	@Tags		deck
//	@Accept		json
//	@Produce	json
//	@Param		matchId		path		int	true	"match id"'
//	@Success	200		{object}	vo.Hand
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/deck/shuffle [put]
func PutShuffle(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	q.ResetDeckForMatchId(ctx, &matchId)

	cardResults, err := q.GetCardsForMatchIdAndState(ctx, queries.GetCardsForMatchIdAndStateParams{
		ID:    &matchId,
		State: queries.CardstateDeck,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	cards := []queries.Matchcard{}
	for _, res := range cardResults {
		cards = append(cards, queries.Matchcard{
			ID:        res.ID,
			Cardid:    res.Cardid,
			Origowner: res.Origowner,
			Currowner: res.Currowner,
			State:     res.State,
		})
	}

	hand := vo.Hand{
		Cards: cards,
	}

	return c.JSON(http.StatusOK, hand)
}
