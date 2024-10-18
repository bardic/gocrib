package player

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/utils"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update kitty with ids
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.HandModifier true "array of ids to add to kitty"
// @Success      200  {object}  queries.Match
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

	err = utils.UpdateGameState(int(m.ID), queries.GamestatePlayState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	m.Gamestate = queries.GamestatePlayState

	return c.JSON(http.StatusOK, m)
}

func updateKitty(details model.HandModifier) (*model.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx := context.Background()

	playerId, err := q.GetCurrentPlayerTurn(ctx, details.MatchId)

	if err != nil {
		return nil, err
	}

	err = q.UpdateKitty(ctx, queries.UpdateKittyParams{
		ID:    playerId,
		Kitty: details.CardIds,
	})

	if err != nil {
		return nil, err
	}

	return utils.PlayCard(details)
}
