package player

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"

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
//	@Param		playerId	path		int	true	"from player id"'
//	@Param		toPlayerId	path		int	true	"to player id"'
// 	@Param		details	body		vo.HandModifier	true	"HandModifier object"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/player/{fromPlayerId}/to/{toPlayerId}/play [put]

func UpdatePlay(c echo.Context) error {
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
	ctx := context.Background()

	// Udpate cards to play state

	for _, cardId := range details.CardIds {
		q.UpdateMatchCardState(ctx, queries.UpdateMatchCardStateParams{
			State:     queries.CardstatePlay,
			Origowner: &fromPlayerId,
			Currowner: &toPlayerId,
			ID:        &cardId,
		})
	}

	//Pass player turn

	orderedPlayers, err := q.GetMatchPlayerOrdered(ctx, &matchId)

	for i, player := range orderedPlayers {
		if *player.ID == toPlayerId {
			playerIndex := i
			nextPlayer := orderedPlayers[(playerIndex+1)%len(orderedPlayers)]
			err := q.UpdateCurrentPlayerTurn(ctx, queries.UpdateCurrentPlayerTurnParams{
				ID:                &matchId,
				Currentplayerturn: nextPlayer.ID,
			})

			if err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := helpers.GetMatch(&matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
