package match

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/player"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new match
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.MatchRequirements true "MatchRequirements"
// @Success      200  {object}  int
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [post]
func NewMatch(c echo.Context) error {
	details := new(model.MatchRequirements)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	p, err := player.NewPlayerQuery(details.AccountId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := q.CreateMatch(ctx, queries.CreateMatchParams{
		Playerids:          []int32{p.ID},
		Privatematch:       false,
		Elorangemin:        0,
		Elorangemax:        0,
		Deckid:             0,
		Cutgamecardid:      0,
		Currentplayerturn:  0,
		Turnpasstimestamps: []pgtype.Timestamptz{},
		Gamestate:          queries.GamestateNewGameState,
		Art:                "default.png",
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.MatchDetailsResponse{
		MatchId:   int(m.ID),
		PlayerId:  int(p.ID),
		GameState: m.Gamestate,
	})
}
