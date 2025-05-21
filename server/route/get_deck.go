package route

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetDeck route
//
//	@Summary			Get decks by matchDd
//	@Description 	Returns vo.Match object when given a matchId
//	@Tags					match
//	@Accept				json
//	@Produce			json
//	@Param				matchId		path			int			true	"Match ID"'
//	@Success			200				{object}	vo.Deck
//	@Failure			400				{object}	error
//	@Failure			500				{object}	error
//	@Router				/match/{matchId}/deck [get]
func (h *Handler) GetDeck(c echo.Context) error {
	p := c.Param("matchId")
	matchID, err := strconv.Atoi(p)
	if err != nil {
		return h.BadParams(c, "", err)
	}

	m, err := h.MatchStore.GetDeck(c, matchID)
	if err != nil {
		return h.InternalError(c, "", err)
	}

	return h.Ok(c, m)
}
