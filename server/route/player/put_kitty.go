package player

import (
	"context"
	"net/http"

	"model"
	"queries"
	conn "server/db"
	"server/utils"

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

	err = utils.UpdateGameState(m.ID, queries.GamestateCutState)

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

	err := q.UpdateKitty(ctx, queries.UpdateKittyParams{
		ID:    details.PlayerId,
		Kitty: details.CardIds,
	})

	if err != nil {
		return nil, err
	}
	err = q.RemoveCardsFromHand(ctx, queries.RemoveCardsFromHandParams{
		ID:   details.PlayerId,
		Hand: details.CardIds,
	})

	if err != nil {
		return nil, err
	}

	m, err := utils.GetMatch(int(details.MatchId))

	if err != nil {
		return nil, err
	}

	return m, nil
}
