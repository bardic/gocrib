package match

import (
	"net/http"

	"server/utils"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Success      200  {object}  []vo.GameMatch
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/matches/open [get]
func GetOpenMatches(c echo.Context) error {
	v, err := utils.GetOpenMatches()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, v)
}
