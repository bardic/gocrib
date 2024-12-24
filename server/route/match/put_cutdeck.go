package match

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// Updates the match with the index of the user selected 'cut' card
//
//	@Summary	Cut deck by index of card selected
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		details	body		vo.CutDeckReq	true	"Deck index that is to become the cut"
//	@Success	200		{object}	int
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/cut [put]
func CutDeck(c echo.Context) error {
	details := new(vo.CutDeckReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cleanId := strings.Trim(details.CutIndex, " ")

	cutIndex, err := strconv.Atoi(cleanId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = q.UpdateMatchCut(ctx, queries.UpdateMatchCutParams{
		ID:            int32(details.MatchId),
		Cutgamecardid: int32(cutIndex),
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.UpdateGameState(details.MatchId, queries.GamestateDealState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, details.MatchId)
}
