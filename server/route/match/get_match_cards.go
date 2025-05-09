package match

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetMatchCardsForMatchID route
// @Summary	Get match card by match id
// @Description
// @Tags		match
// @Accept		json
// @Produce	json
// @Param		matchId	path		int	true	"match id"'
// @Success	200	{object}	[]queries.GetCardsForMatchIdRow
// @Failure	404	{object}	error
// @Failure	422	{object}	error
// @Router		/match/{matchId}/cards [get]
func (h *Handler) GetMatchCardsForMatchID(c echo.Context) error {
	matchID, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	cards, err := h.MatchStore.GetCardsForMatchID(c, matchID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, cards)
}
