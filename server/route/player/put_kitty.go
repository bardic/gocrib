package player

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"

	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	Update kitty with ids
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId		path		int	true	"match id"'
//	@Param		playerId	path		int	true	"from player id"'
//	@Param		toPlayerId	path		int	true	"to player id"'
//	@Param		details	body		vo.HandModifier	true	"array of ids to add to kitty"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/player/{fromPlayerId}/to/{toPlayerId}/kitty [put]
func UpdateKitty(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	fromPlayerId, err := strconv.Atoi(c.Param("playerId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	toPlayerId, err := strconv.Atoi(c.Param("toPlayerId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	details := &vo.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, cardId := range details.CardIds {
		q.UpdateMatchCardState(ctx, queries.UpdateMatchCardStateParams{
			ID:        &cardId,
			State:     queries.CardstateKitty,
			Origowner: &fromPlayerId,
			Currowner: &toPlayerId,
		})
	}

	err = q.MarkPlayerReady(ctx, queries.MarkPlayerReadyParams{
		Isready: true,
		ID:      &fromPlayerId,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := helpers.GetMatch(&matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	allReady := true
	for _, player := range m.Players {
		if !player.Isready {
			allReady = false
			break
		}
	}

	if allReady {
		err = q.UpdateMatchState(ctx, queries.UpdateMatchStateParams{
			Gamestate: queries.GamestateCut,
			ID:        &matchId,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, m)
}
