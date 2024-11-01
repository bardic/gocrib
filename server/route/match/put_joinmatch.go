package match

import (
	"context"
	"net/http"
	"time"

	"queries"
	conn "server/db"
	"server/route/player"
	"vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Join match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body vo.JoinMatchReq true "match Object to update"
// @Success      200  {object}  vo.MatchDetailsResponse
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/join [put]
func JoinMatch(c echo.Context) error {
	details := new(vo.JoinMatchReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := player.NewPlayerQuery(details.AccountId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = q.JoinMatch(ctx, queries.JoinMatchParams{
		Matchid:  details.MatchId,
		Playerid: p.ID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	match, err := q.UpdateGameState(ctx, queries.UpdateGameStateParams{
		ID:        details.MatchId,
		Gamestate: queries.GamestateJoinGameState,
	})	

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, match)
}
