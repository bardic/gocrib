package match

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

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
func (h *MatchHandler) GetPlayerIdForMatchAndAccount(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	accountId, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	playerData, err := h.PlayerStore.GetPlayerByMatchAndAccountId(c, queries.GetPlayerByMatchAndAccountIdParams{
		Matchid:   &matchId,
		Accountid: &accountId,
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
