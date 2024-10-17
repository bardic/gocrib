package gameplaycard

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/gocrib/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get cards by ids
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param        ids    query     string  true  "csv of ids"'
// @Success      200  {object}  []model.GameplayCard{}
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/gameplaycards/ [get]
func GetGameplayCards(c echo.Context) error {
	ids := c.Request().URL.Query().Get("ids")
	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM gameplaycards NATURAL JOIN cards WHERE gameplaycards.id IN ("+ids+")")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	v := []model.GameplayCard{}

	for rows.Next() {
		var card model.GameplayCard

		err := rows.Scan(&card.Id, &card.CardId, &card.OrigOwner, &card.CurrOwner, &card.State, &card.Value, &card.Suit, &card.Art)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		v = append(v, card)
	}

	return c.JSON(http.StatusOK, v)
}
