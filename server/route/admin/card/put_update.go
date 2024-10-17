package card

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/model"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update card by id
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param details body model.Card true "card Object to save"
// @Success      200  {object}  model.Card
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /admin/card/ [put]
func UpdateCard(c echo.Context) error {
	details := new(model.Card)
	if err := c.Bind(details); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	args := parseCard(*details)

	query := "UPDATE cards SET value=@value, suit=@suit, art=@art where id = @id"

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Card updated")
}
