package match

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"

	"github.com/labstack/echo/v4"
)

// Updates the match with the index of the user selected 'cut' card
//
//	@Summary	Update the matches current palyer
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"match id"
//	@Param		playerId	path		int	true	"playerId id"
//	@Success	200		{object}	int
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/currentPlayer/{playerId} [put]
func (h *MatchHandler) UpdateCurrentPLayer(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	playerId, err := strconv.Atoi(c.Param("playerId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = q.UpateMatchCurrentPlayerTurn(ctx, queries.UpateMatchCurrentPlayerTurnParams{
		ID:                &matchId,
		Currentplayerturn: &playerId,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := helpers.GetMatch(&matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}
