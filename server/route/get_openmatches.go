package route

import (
	"github.com/bardic/gocrib/vo"
	"github.com/labstack/echo/v4"
)

// GetOpenMatches route
//
//	@Summary			Get list of open matches
//	@Description
//	@Tags					match
//	@Accept				json
//	@Produce			json
//	@Success			200			{object}	[]vo.Match
//	@Failure			400			{object}	error
//	@Failure			500			{object}	error
//	@Router				/open 	[get]
func (h *Handler) GetOpenMatches(c echo.Context) error {
	m, err := h.MatchStore.GetOpenMatches(c)
	if err != nil {
		return h.InternalError(c, "failed to get open matches", err)
	}

	if m == nil {
		m = []*vo.Match{}
	}

	return h.Ok(c, m)
}
