package route

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/labstack/echo/v4"
)

// CutDeck route
//
//	@Summary	Cut deck by index of card selected
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"match id"
//	@Param		cutIndex	path		int	true	"cut id"
//	@Success	200		{object}	int
//	@Failure	400		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/cut/{cutId} [put]
func (h *Handler) CutDeck(c echo.Context) error {
	matchID, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cutIndex, err := strconv.Atoi(c.Param("cutIndex"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = h.MatchStore.UpdateMatchCut(c, queries.UpdateMatchCutParams{
		ID:            matchID,
		Cutgamecardid: cutIndex,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = UpdateGameState(matchID, queries.GamestatePlay)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := GetMatch(matchID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
