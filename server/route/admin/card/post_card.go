package card

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/gocrib/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new card
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param details body model.Card true "card Object to save"
// @Success      200  {object}  model.Card
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /admin/card/ [post]
func NewCard(c echo.Context) error {
	details := new(model.Card)
	if err := c.Bind(details); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	args := parseCard(*details)
	query := "INSERT INTO cards(value, suit, art) VALUES (@value, @suit, @art)"

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Success")
}
