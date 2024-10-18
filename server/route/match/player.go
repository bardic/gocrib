package match

import (
	"net/http"

	"github.com/bardic/gocrib/server/utils"

	"github.com/bardic/gocrib/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update play with ids
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.HandModifier true "array of ids to add to play"
// @Success      200  {object}  queries.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/play [put]
func UpdatePlay(c echo.Context) error {
	details := &model.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := utils.UpdatePlay(*details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = utils.UpdateGameState(m.Id, model.OpponentState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
