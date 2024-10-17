package match

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbage/server/utils"
	"github.com/bardic/gocrib/model"
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
// @Success      200  {object}  model.MatchDetailsResponse
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/join [put]
func JoinMatch(c echo.Context) error {
	details := new(model.JoinMatchReq)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := utils.GetPlayerById(details.PlayerId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	details.PlayerId = p.Id
	m, err := updatePlayersInMatch(*details)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// // temp force players to ready
	// _, err = player.ReadyPlayerById(c, p.Id)

	// if err != nil {
	// 	return err
	// }

	// m, err = utils.GetMatch(m.Id)

	// if err != nil {
	// 	return err
	// }

	// //check if players are ready
	// rdy := utils.PlayersReady(m.Players)

	// gameState := model.JoinGameState
	// if rdy {
	// 	game.Deal(m)
	// }

	// utils.UpdateMatchState(details.MatchId, gameState)

	return c.JSON(http.StatusOK, model.MatchDetailsResponse{
		MatchId:   m.Id,
		GameState: model.JoinGameState,
	})
}

func updatePlayersInMatch(req model.JoinMatchReq) (*model.GameMatch, error) {
	args := pgx.NamedArgs{
		"matchId":  req.MatchId,
		"playerId": req.PlayerId,
	}

	query := `UPDATE match SET
				playerIds=ARRAY_APPEND(playerIds, @playerId)
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

	return m, nil
}
