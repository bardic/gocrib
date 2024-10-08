package match

import (
	"context"
	"encoding/json"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbage/server/route/game"
	"github.com/bardic/cribbage/server/route/player"
	"github.com/bardic/cribbage/server/utils"
	"github.com/bardic/cribbagev2/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Join match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.JoinMatchReq true "match Object to update"
// @Success      200  {object}  model.GameMatch
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/join [put]
func JoinMatch(c echo.Context) error {
	details := new(model.JoinMatchReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := player.NewPlayerQuery(details.RequesterId)

	if err != nil {
		return err
	}

	details.PlayerId = p.Id
	m, err := updatePlayersInMatch(*details)
	if err != nil {
		return err
	}

	b, err := json.Marshal(m)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func updatePlayersInMatch(req model.JoinMatchReq) (*model.GameMatch, error) {
	args := pgx.NamedArgs{
		"matchId":     req.MatchId,
		"requesterId": req.RequesterId,
		"playerId":    req.PlayerId,
		"gameState":   model.WaitingState,
	}

	query := `UPDATE match SET
				playerIds=ARRAY_APPEND(playerIds, @playerId),
				gameState=@gameState
			WHERE id=@matchId`

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return nil, err
	}

	m, err := utils.GetMatch(req.MatchId)

	if err != nil {
		return nil, err
	}

	//TODO: Not here
	_, err = game.Deal(m)

	if err != nil {
		return nil, err
	}

	return &m, nil
}
