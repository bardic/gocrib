package match

import (
	"net/http"
	"strconv"

	"server/utils"
	"vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by id"'
// @Success      200  {object}  vo.MatchDetailsResponse
// @Failure      404  {object}  error
// @Failure      422  {object}  error
// @Router       /player/match/ [get]
func GetMatchState(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	m, err := utils.GetMatch(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, vo.MatchDetailsResponse{
		MatchId:   int(m.ID),
		GameState: m.Gamestate,
	})
}
