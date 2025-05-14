package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/labstack/echo/v4"
)

// Updates the match with the index of the user selected 'cut' card
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
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/play [put]
func (h *Handler) Play(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cutIndex, err := strconv.Atoi(c.Param("cutIndex"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = h.MatchStore.UpdateMatchCut(c, queries.UpdateMatchCutParams{
		ID:            &matchId,
		Cutgamecardid: &cutIndex,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = UpdateGameState(&matchId, queries.GamestatePlay)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := GetMatch(&matchId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
