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
//	@Summary	Cut deck by index of card selected
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"match id"
//	@Param		cutIndex	path		int	true	"cut id"
//	@Success	200		{object}	int
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/cut/{cutId} [put]
func CutDeck(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	cutIndex, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = q.UpdateMatchCut(ctx, queries.UpdateMatchCutParams{
		ID:            int32(matchId),
		Cutgamecardid: int32(cutIndex),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.UpdateGameState(int32(matchId), queries.GamestateDeal)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, matchId)
}
