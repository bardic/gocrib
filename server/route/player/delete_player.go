package player

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get player by barcode
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [delete]
func DeletePlayer(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}
