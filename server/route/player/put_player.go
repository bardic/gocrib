package player

import (
	"net/http"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/server/utils"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update player by barcode
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body model.Player true "player Object to save"
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [put]
func UpdatePlayer(c echo.Context) error {
	details := new(model.Player)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := utils.UpdatePlayerById(*details)

	if err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, p)
}
