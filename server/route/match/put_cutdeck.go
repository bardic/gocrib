package match

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"queries"
	conn "server/db"
	"server/utils"
	"vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Cut deck by index
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body vo.CutDeckReq true "Deck index that is to become the cut"
// @Success      200  {object}  int
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/cut [put]
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

	err = utils.UpdateGameState(details.MatchId, queries.GamestateDealState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, details.MatchId)
}
