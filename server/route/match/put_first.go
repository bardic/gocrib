package match

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"

	conn "github.com/bardic/gocrib/server/db"

	"github.com/labstack/echo/v4"
)

// Return match details if the player is able to join the match
//
//	@Summary	Join match by id
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		matchId	path		int	true	"match id"'
//	@Success	200		{object}	vo.MatchDetailsResponse
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/{matchId}/determinefirst/ [put]
func DetermineFirst(c echo.Context) error {
	matchId, err := strconv.Atoi(c.Param("matchId"))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	match, err := OnDetermineFirst(matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, match)
}

func OnDetermineFirst(matchId int) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m, err := q.GetMatchById(ctx, &matchId)

	if err != nil {
		return nil, err
	}

	var match *vo.GameMatch
	err = json.Unmarshal(m, &match)

	if err != nil {
		return nil, err
	}

	players, err := q.GetMatchPlayersByMatchId(ctx, &matchId)

	if err != nil {
		return nil, err
	}

	turnOrder := 1
	for _, player := range players {
		q.UpdatePlayerTurnOrder(ctx, queries.UpdatePlayerTurnOrderParams{
			Turnorder: &turnOrder,
			Matchid:   &matchId,
			Playerid:  player.ID,
		})

		if turnOrder == 1 {

			err = q.UpdateCurrentPlayerTurn(ctx, queries.UpdateCurrentPlayerTurnParams{
				ID:                &matchId,
				Currentplayerturn: player.ID,
			})

			if err != nil {
				return nil, err
			}

			err = q.UpdateDealerForMatch(ctx, queries.UpdateDealerForMatchParams{
				ID:       &matchId,
				Dealerid: player.ID,
			})

			if err != nil {
				return nil, err
			}
		}

		turnOrder++
	}

	updatedMatch, err := q.UpdateGameState(ctx, queries.UpdateGameStateParams{
		ID:        &matchId,
		Gamestate: queries.GamestateDeal,
	})

	p := []vo.GamePlayer{}
	for _, player := range players {
		p = append(p, vo.GamePlayer{
			Player: queries.Player{
				ID:        player.ID,
				Accountid: player.Accountid,

				Isready: player.Isready,
			},
			Hand:  []queries.Matchcard{},
			Play:  []queries.Matchcard{},
			Kitty: []queries.Matchcard{},
		})
	}

	gameMatch := vo.GameMatch{
		Match:   updatedMatch,
		Players: p,
	}

	if err != nil {
		return nil, err
	}

	return &gameMatch, nil
}
