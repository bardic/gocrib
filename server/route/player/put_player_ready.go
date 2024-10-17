package player

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/model"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/utils"
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
// @Success      200  {object}  bool
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/ready [put]
func PlayerReady(c echo.Context) error {
	details := new(int)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := ReadyPlayerById(c, *details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	matchId, err := utils.GetMatchForPlayerId(*details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m, err := utils.GetMatch(matchId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if utils.PlayersReady(m.Players) {
		utils.Deal(m)
		utils.UpdateGameState(matchId, model.DiscardState)
	}

	return c.JSON(http.StatusOK, nil)
}

func ReadyPlayerById(c echo.Context, playerId int) (bool, error) {
	args := pgx.NamedArgs{
		"id":      playerId,
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
		return false, err
	}

	return true, nil
}
