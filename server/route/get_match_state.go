package route

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetMatchState route
//
//	@Summary			Get state by matchId
//	@Description
//	@Tags					match
//	@Accept				json
//	@Produce			json
//	@Param				matchId				path		int	true	"Match ID"'
//	@Success			200	{object}	string
//	@Failure			400	{object}	error
//	@Failure			500	{object}	error
//	@Router				/match/{matchId}/state [get]
func (h *Handler) GetMatchState(c echo.Context) error {
	p := c.Param("matchId")
	matchID, err := strconv.Atoi(p)
	if err != nil {
		return h.BadParams(c, "", err)
	}

	m, err := h.MatchStore.GetMatchState(c, matchID)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	return h.Ok(c, m)
}
