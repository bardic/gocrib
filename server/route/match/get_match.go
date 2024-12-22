package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/server/route/helpers"
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

// Returns a match by id
//
//	@Summary	Get match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce		json
//	@Param		id	path		int	true	"Match ID"'
//	@Success	200	{object}	queries.Match
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/match/{id} [get]
func GetMatch(c echo.Context) error {
	p := c.Param("id")
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
