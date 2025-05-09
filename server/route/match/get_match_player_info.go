package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// GetPlayerIDForMatchAndAccount route
// @Summary	Get match card by match id
// @Description
// @Tags		match
// @Accept		json
// @Produce	json
// @Param		matchId	path		int	true	"match id"'
// @Param		accountId	path		int	true	"account id"'
// @Success	200	{object}	vo.GamePlayer
// @Failure	404	{object}	error
// @Failure	422	{object}	error
// @Router		/match/{matchId}/account/{accountId} [get]
func (h *Handler) GetPlayerIDForMatchAndAccount(c echo.Context) error {
	matchID, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	playerData, err := h.PlayerStore.GetPlayerByMatchAndAccountID(c, queries.GetPlayerByMatchAndAccountIdParams{
		Matchid:   &matchID,
		Accountid: &accountID,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	gamePlayer := vo.GamePlayer{
		Player: &queries.Player{
			ID:        playerData.ID,
			Accountid: playerData.Accountid,
			Score:     playerData.Score,
			Isready:   playerData.Isready,
			Art:       playerData.Art,
		},
	}

	return c.JSON(http.StatusOK, gamePlayer)
}
