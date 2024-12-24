package player

import (
	"context"
	"net/http"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"
	"github.com/bardic/gocrib/server/route/helpers"
	"github.com/bardic/gocrib/vo"

	"github.com/labstack/echo/v4"
)

// Create godoc
//
//	@Summary	Update kitty with ids
//	@Description
//	@Tags		match
//	@Accept		json
//	@Produce	json
//	@Param		details	body		vo.HandModifier	true	"array of ids to add to kitty"
//	@Success	200		{object}	queries.Match
//	@Failure	400		{object}	error
//	@Failure	404		{object}	error
//	@Failure	500		{object}	error
//	@Router		/match/player//kitty [put]
func UpdateKitty(c echo.Context) error {
	details := &vo.HandModifier{}
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m, err := updateKitty(*details)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = helpers.UpdateGameState(m.ID, queries.GamestateCutState)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}

func updateKitty(details vo.HandModifier) (*vo.GameMatch, error) {
	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

	m, err := helpers.GetMatch(int(details.MatchId))

	if err != nil {
		return nil, err
	}

	return m, nil
}
