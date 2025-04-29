package player

import (
	"net/http"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/server/route/match"
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
func (h *PlayerHandler) PlayerReady(c echo.Context) error {
	pReady := new(vo.PlayerReady)
	if err := c.Bind(pReady); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.readyPlayerById(c, pReady.PlayerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := match.GetMatch(pReady.MatchId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if arePlayersReady(m.Players) {
		// TODO
		// /controller.Deal(match)
		match.UpdateGameState(m.ID, queries.GamestateDiscard)
	}

	return c.JSON(http.StatusOK, nil)
}

func arePlayersReady(players []*vo.GamePlayer) bool {
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

func (h *PlayerHandler) readyPlayerById(ctx echo.Context, playerId *int) (bool, error) {
	err := h.PlayerStore.UpdatePlayerReady(ctx, queries.UpdatePlayerReadyParams{
		ID:      playerId,
		Isready: true,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
