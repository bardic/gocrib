package card

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update card by id
// @Description
// @Tags         cards
// @Accept       json
// @Produce      json
// @Param details body queries.Card true "card Object to save"
// @Success      200  {object}  queries.Card
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /admin/card/ [put]
func UpdateCard(c echo.Context) error {

	return c.JSON(http.StatusOK, "Card updated")
}
