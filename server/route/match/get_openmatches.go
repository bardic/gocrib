package match

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetOpenMatches route
//
//	@Summary	Get match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	[]vo.GameMatch
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/open [get]
func (h *Handler) GetOpenMatches(c echo.Context) error {
	v, err := GetOpenMatches()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, v)
}
