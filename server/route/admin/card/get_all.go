package card

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get all cards
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Success      200  {object}  []queries.Card{}
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/allcards/ [get]
func GetAllCards(c echo.Context) error {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	cards, err := q.GetCards(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cards)
}
