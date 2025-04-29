package match

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Returns a match by id
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
//	@Router		/match/{matchId} [get]
func (h *MatchHandler) GetMatch(c echo.Context) error {
	p := c.Param("matchId")
	matchId, err := strconv.Atoi(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	m, err := h.MatchStore.GetMatch(c, &matchId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
