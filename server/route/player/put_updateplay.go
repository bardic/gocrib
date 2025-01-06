package player

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/server/controller"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"

	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// Updates a matches cards state
//
//	@Summary	Update play
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId		path		int	true	"match id"'
//	@Param		playerId	path		int	true	"player id"'
//	@Param		details	body		vo.HandModifier	true	"HandModifier object"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/player/{playerId}/play [put]
func UpdatePlay(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	playerId, err := strconv.Atoi(c.Param("playerId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	details := &vo.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := controller.UpdatePlay(matchId, playerId, *details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.UpdateGameState(m.ID, queries.GamestatePassTurn)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = q.PassTurn(ctx, queries.PassTurnParams{
		Matchid:  int32(matchId),
		Playerid: int32(playerId),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}