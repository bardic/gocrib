package route

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// Play route
//
//	@Summary 			Update the state of several cards
//	@Description
//	@Tags					match
//	@Accept				json
//	@Produce			json
//	@Param				matchId		path				int	true	"match id"
//	@Param				cutIndex	path				int	true	"cut id"
//	@Success			200				{object}		int
//	@Failure			400				{object}		error
//	@Failure			500				{object}		error
//	@Router				/match/{matchId}/play [put]
func (h *Handler) Play(c echo.Context) error {
	matchID, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return h.InternalError(c, "", err)
	}

	cutIndex, err := strconv.Atoi(c.Param("cutIndex"))
	if err != nil {
		return h.InternalError(c, "", err)
	}

	err = h.MatchStore.UpdateMatchCut(c, matchID, cutIndex)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	err = h.MatchStore.SetMatchState(c, matchID, "Play")
	if err != nil {
		return h.InternalError(c, "", err)
	}

	m, err := h.MatchStore.GetMatch(c, matchID)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	return h.Ok(c, m)
}
