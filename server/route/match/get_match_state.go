package match

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetMatch route
//
//	@Summary	Get match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce		json
//	@Param		matchId	path		int	true	"Match ID"'
//	@Success	200	{object}	queries.Match
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/match/{matchId}/state [get]
func (h *Handler) GetMatchState(c echo.Context) error {
	p := c.Param("matchId")
	matchID, err := strconv.Atoi(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	m, err := h.MatchStore.GetMatchState(c, &matchID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
