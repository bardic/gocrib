package match

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"queries"
	conn "server/db"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get match cards by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by id"'
// @Success      200  {object}  []queries.GetMatchCardsRow
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/match/cards/ [get]
func GetGameCardsForMatch(c echo.Context) error {
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

	//cards, err := q.GetGameCardsForMatch(ctx, int32(id))
	cards, err := q.GetMatchCards(ctx, int32(id))

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, cards)
}
