package route

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
//	@Success	200	{object}	vo.Match
//	@Failure	400	{object}	error
//	@Failure	500	{object}	error
//	@Router		/match/{matchId} [get]
func (h *Handler) GetMatch(c echo.Context) error {
	p := c.Param("matchId")
	matchID, err := strconv.Atoi(p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	m, err := h.MatchStore.GetMatch(c, matchID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
