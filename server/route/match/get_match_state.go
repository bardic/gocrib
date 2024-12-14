package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/server/route/helpers"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// Takes match id and returns current state
//
//	@Summary	Get match state by match id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		id	query		string	true	"match id"'
//	@Success	200	{object}	vo.MatchDetailsResponse
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/player/match/ [get]
func GetMatchState(c echo.Context) error {
	p := c.Request().URL.Query().Get("id")
	id, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	m, err := helpers.GetMatch(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, vo.MatchDetailsResponse{
		MatchId:   m.ID,
		GameState: m.Gamestate,
	})
}
