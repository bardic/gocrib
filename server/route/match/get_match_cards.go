package match

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"queries"

	conn "github.com/bardic/gocrib/server/db"

	"github.com/labstack/echo/v4"
)

// Takes a match id and returns the cards for that match
//
//	@Summary	Get match cards by match id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		id	query		string	true	"match id"'
//	@Success	200	{object}	[]queries.GetMatchCardsRow
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/match/cards/ [get]
func GetMatchCardsForMatchId(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cards, err := q.GetMatchCards(ctx, int32(id))

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, cards)
}
