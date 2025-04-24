package match

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	conn "github.com/bardic/gocrib/server/db"

	"github.com/labstack/echo/v4"
)

// Takes a match and account id and returns the player
//
//	@Summary	Get match card by match id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"match id"'
//	@Param		accountId	path		int	true	"account id"'
//	@Success	200	{object}	vo.GamePlayer
//	@Failure	404	{object}	error
//	@Failure	422	{object}	error
//	@Router		/match/{matchId}/account/{accountId} [get]
func GetPlayerIdForMatchAndAccount(c echo.Context) error {
	p := c.Param("matchId")
	matchId, err := strconv.Atoi(p)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	p2 := c.Param("accountId")
	accountId, err := strconv.Atoi(p2)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playerData, err := q.GetPlayerByMatchAndAccountId(ctx, queries.GetPlayerByMatchAndAccountIdParams{
		Matchid:   &matchId,
		Accountid: &accountId,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	gamePlayer := vo.GamePlayer{
		Player: queries.Player{
			ID:        playerData.ID,
			Accountid: playerData.Accountid,
			Score:     playerData.Score,
			Isready:   playerData.Isready,
			Art:       playerData.Art,
		},
	}

	return c.JSON(http.StatusOK, gamePlayer)
}
