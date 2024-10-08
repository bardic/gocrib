package deck

import (
	"context"
	"net/http"
	"strconv"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbagev2/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get deck by id
// @Description
// @Tags         deck
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for deck by id"'
// @Success      200  {object}  model.GameDeck
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/match/deck/ [get]
func GetDeck(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return err
	}

	deck, err := GetDeckById(id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, deck)
}

func GetDeckById(id int) (model.GameDeck, error) {
	db := conn.Pool()
	defer db.Close()
	var deckId int
	var cards []model.GameplayCard
	err := db.QueryRow(context.Background(), "SELECT * FROM deck WHERE id=$1", id).Scan(&deckId, &cards)

	if err != nil {
		return model.GameDeck{}, err
	}

	deck := model.GameDeck{
		Id:    deckId,
		Cards: cards,
	}

	return deck, nil
}
