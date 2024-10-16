package player

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
// @Summary      Update player to mark as ready
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body int true "player id to update"
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/ready [put]
func PlayerReady(c echo.Context) error {
	details := new(int)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := pgx.NamedArgs{
		"id":      details,
		"isReady": true,
	}

	query := `UPDATE player SET 
		isReady = @isReady
	where 
		id = @id`

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	args = pgx.NamedArgs{
		"id": details,
	}

	query = `SELECT id, playerids from match WHERE @id = ANY(playerids)`

	var matchId int
	var playerIds []int
	err = db.QueryRow(
		context.Background(),
		query,
		args).Scan(&matchId, &playerIds)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	allReady := true
	for _, playerId := range playerIds {
		args = pgx.NamedArgs{
			"id": playerId,
		}
		query = `SELECT isReady from player WHERE id = @id`

		var isReady bool
		err = db.QueryRow(
			context.Background(),
			query,
			args).Scan(&isReady)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		if !isReady {
			allReady = false
			break
		}

	}

	if allReady {
		utils.UpdateGameState(matchId, model.CutState)
	}

	return c.JSON(http.StatusOK, nil)
}
