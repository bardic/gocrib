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
// @Summary      Update kitty with ids
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.HandModifier true "array of ids to add to kitty"
// @Success      200  {object}  model.GameMatch
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/kitty [put]
func UpdateKitty(c echo.Context) error {
	details := &model.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := updateKitty(*details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = utils.UpdateGameState(m.Id, model.PlayState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m.GameState = model.PlayState

	return c.JSON(http.StatusOK, m)
}

func updateKitty(details model.HandModifier) (*model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()

	args := pgx.NamedArgs{
		"matchId": details.MatchId,
	}

	q := "SELECT currentplayerturn FROM match WHERE id = @matchId"

	var dealerPlayerId int
	err := db.QueryRow(
		context.Background(),
		q,
		args).Scan(&dealerPlayerId)

	if err != nil {
		return nil, err
	}

	args = pgx.NamedArgs{
		"dealerId": 1,
		"kitty":    details.CardIds,
	}

	q = "UPDATE player SET kitty = kitty + @kitty where id = @dealerId"

	_, err = db.Exec(
		context.Background(),
		q,
		args)

	if err != nil {
		return nil, err
	}

	return utils.PlayCard(details)
}
