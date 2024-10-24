package match

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/utils"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Cut deck by index
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.CutDeckReq true "Deck index that is to become the cut"
// @Success      200  {object}  int
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/cut [put]
func CutDeck(c echo.Context) error {
	details := new(model.CutDeckReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

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

	err = utils.UpdateGameState(details.MatchId, queries.GamestatePlayState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, details.MatchId)
}
