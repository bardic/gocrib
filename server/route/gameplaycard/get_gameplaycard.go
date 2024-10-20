package gameplaycard

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get cards by ids
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        ids    query     string  true  "csv of ids"'
// @Success      200  {object}  []queries.Card{}
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/gameplaycards/ [get]
func GetGameplayCards(c echo.Context) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	strIds := strings.Split(c.Request().URL.Query().Get("ids"), ",")

	ids := []int32{}
	for _, i := range strIds {
		_id, err := strconv.Atoi(i)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		ids = append(ids, int32(_id))
	}

	cards, err := q.GetMatchCards(ctx, ids)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cards)
}
