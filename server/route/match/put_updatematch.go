package match

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbage/server/utils"
	"github.com/bardic/gocrib/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update match by id
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.GameMatch true "match Object to update"
// @Success      200  {object}  model.GameMatch
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/ [put]
func UpdateMatch(c echo.Context) error {
	details := new(model.GameMatch)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := updateMatch(*details); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

func updateMatch(match model.GameMatch) error {
	args := utils.ParseMatch(match)
	query := `UPDATE match SET
				playerIds = @playerIds,
				creationDate = @creationDate,
				privateMatch = @privateMatch,
				eloRangeMin = @eloRangeMin,
				eloRangeMax = @eloRangeMax,
				deckId = @deckId,
				cutGameCardId = @cutGameCardId,
				currentPlayerTurn = @currentPlayerTurn,
				turnPassTimestamps = @turnPassTimestamps,
				gameState= @gameState,
				art = @art
			WHERE id=@id`

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return err
	}

	return nil
}
