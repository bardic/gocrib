package match

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"
	"github.com/bardic/gocrib/server/route/player"

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
//	@Param		accountId	path		int	true	"account id"'
//	@Success	200		{object}	vo.MatchDetailsResponse
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/join/{accountId} [put]
func JoinMatch(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accountId, err := strconv.Atoi(c.Param("accountId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	player, err := player.NewPlayerQuery(&matchId, &accountId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = q.JoinMatch(ctx, queries.JoinMatchParams{
		Matchid:  &matchId,
		Playerid: player.ID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	_, err = q.UpdateGameState(ctx, queries.UpdateGameStateParams{
		ID:        &matchId,
		Gamestate: queries.GamestateDetermine,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	_, err = helpers.GetMatch(&matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	_, err = OnDetermineFirst(matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := OnDeal(matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, match)
}
