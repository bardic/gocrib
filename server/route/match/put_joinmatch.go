package match

import (
	"net/http"

	"queries"
	"server/utils"
	"vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Join match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body vo.JoinMatchReq true "match Object to update"
// @Success      200  {object}  vo.MatchDetailsResponse
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/join [put]
func JoinMatch(c echo.Context) error {
	details := new(vo.JoinMatchReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := utils.GetPlayerById(details.PlayerId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, vo.MatchDetailsResponse{
		MatchId:   int(details.MatchId),
		PlayerId:  int(p.ID),
		GameState: queries.GamestateJoinGameState,
	})
}
