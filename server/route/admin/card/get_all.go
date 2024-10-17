package card

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/model"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get all cards
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Success      200  {object}  []model.GameplayCard{}
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/allcards/ [get]
func GetAllCards(c echo.Context) error {
	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM cards")

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	v := []model.Card{}

	for rows.Next() {
		var card model.Card

		err := rows.Scan(&card.Id, &card.Value, &card.Suit, &card.Art)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		v = append(v, card)
	}

	if len(v) == 0 {
		return c.JSON(http.StatusNotFound, "No cards found. Something is very wrong")
	}

	return c.JSON(http.StatusOK, v)
}
