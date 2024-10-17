package match

import (
	"net/http"

	"github.com/bardic/cribbage/server/utils"
	"github.com/bardic/gocrib/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.GameMatch true "match Object to update"
// @Success      200  {object}  model.GameMatch
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [put]
func UpdateMatch(c echo.Context) error {
	details := new(model.GameMatch)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := utils.UpdateMatch(*details); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "meow")
}
