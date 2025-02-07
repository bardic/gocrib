package match

import (
	"context"
	"encoding/json"
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
//	@Success	200	{object}	[]queries.GetMatchCardsRow
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

	playerData, err := q.GetPlayerByAccountAndMatchIdJSON(ctx, queries.GetPlayerByAccountAndMatchIdJSONParams{
		Matchid:   &matchId,
		Accountid: &accountId,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	var gamePlayer *vo.GamePlayer
	err = json.Unmarshal(playerData, &gamePlayer)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, gamePlayer)
}
