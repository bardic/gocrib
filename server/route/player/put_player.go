package player

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbage/server/utils"
	"github.com/bardic/gocrib/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Update player by barcode
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body model.Player true "player Object to save"
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [put]
func UpdatePlayer(c echo.Context) error {
	details := new(model.Player)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := UpdatePlayerById(*details)

	if err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, p)
}

func UpdatePlayerById(player model.Player) (model.Player, error) {
	args := utils.ParsePlayer(player)

	query := `UPDATE player SET 
		hand = @hand, 
		play = @play, 
		kitty = @kitty, 
		score = @score, 
		isReady = @isReady,
		art = @art 
	where 
		id = @id`

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return model.Player{}, err
	}

	return player, nil
}
