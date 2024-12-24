package match

import (
	"context"
	"net/http"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/player"
	"github.com/bardic/gocrib/vo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	Create new match
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		details	body		vo.MatchRequirements	true	"MatchRequirements"
//	@Success	200		{object}	int
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/ [post]
func NewMatch(c echo.Context) error {
	details := new(vo.MatchRequirements)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := player.NewPlayerQuery(details.AccountId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	deck, err := q.CreateDeck(ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := q.CreateMatch(ctx, queries.CreateMatchParams{
		Privatematch:       false,
		Elorangemin:        0,
		Elorangemax:        0,
		Deckid:             deck.ID,
		Cutgamecardid:      0,
		Currentplayerturn:  p.ID,
		Turnpasstimestamps: []pgtype.Timestamptz{},
		Gamestate:          queries.GamestateNewGameState,
		Art:                "default.png",
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = q.JoinMatch(ctx, queries.JoinMatchParams{
		Matchid:  m.ID,
		Playerid: p.ID,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, vo.MatchDetailsResponse{
		MatchId:   m.ID,
		PlayerId:  p.ID,
		GameState: m.Gamestate,
	})
}
