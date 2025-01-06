package player

import (
	"context"
	"net/http"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/server/controller"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// PReady struct

// Update player by id to be ready. Returns true if all players are ready
//
//	@Summary	Update player by id to be ready. Returns true if all players are ready
//	@Description
//	@Tags		players
//	@Accept		json
//	@Produce	json
//	@Param		pReady	body		vo.PlayerReady	true	"player id to update"
//	@Success	200		{object}	bool
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/player/ready [put]
func PlayerReady(c echo.Context) error {
	pReady := new(vo.PlayerReady)
	if err := c.Bind(pReady); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := readyPlayerById(int(pReady.PlayerId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := helpers.GetMatch(int(pReady.MatchId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if arePlayersReady(match.Players) {
		controller.Deal(match)
		helpers.UpdateGameState(match.ID, queries.GamestateDiscard)
	}

	return c.JSON(http.StatusOK, nil)
}

func arePlayersReady(players []*queries.Player) bool {
	ready := true

	// if len(players) < 2 {
	// 	return false
	// }

	// for _, p := range players {
	// 	if !p.Isready {
	// 		ready = false
	// 	}
	// }

	return ready
}

func readyPlayerById(playerId int) (bool, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := q.UpdatePlayerReady(ctx, queries.UpdatePlayerReadyParams{
		ID:      int32(playerId),
		Isready: true,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func isPlayerReady(playerId int) (bool, error) {
	// db := conn.Pool()
	// defer db.Close()
	// q := queries.New(db)

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// p, err := q.GetPlayerById(ctx, int32(playerId))

	// if err != nil {
	// 	return false, err
	// }

	// return p.Isready, nil

	return true, nil
}
