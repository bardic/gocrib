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

// Return match details if the player is able to join the match
//
//	@Summary	Join match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"match id"'
//	@Success	200		{object}	vo.MatchDetailsResponse
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/pass [put]
func (h *MatchHandler) Pass(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := q.GetMatchById(ctx, &matchId)

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	var match *vo.GameMatch
	err = json.Unmarshal(m, &match)

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	nextPlayer, err := q.GetNextPlayerInTurnOrder(ctx, queries.GetNextPlayerInTurnOrderParams{
		Matchid: &matchId,
		Column2: &match.Currentplayerturn,
	})

	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	err = q.UpateMatchCurrentPlayerTurn(ctx, queries.UpateMatchCurrentPlayerTurnParams{
		ID:                &matchId,
		Currentplayerturn: nextPlayer.ID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = q.UpdateDealerForMatch(ctx, queries.UpdateDealerForMatchParams{
		ID:       &matchId,
		Dealerid: nextPlayer.ID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	updatedMatch, err := q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{
		ID:        &matchId,
		Gamestate: queries.GamestateDeal,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, updatedMatch)
}
