package player

import (
	"context"
	"net/http"
	"strconv"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbagev2/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get player by barcode
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [get]
func GetPlayer(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	p1Id, err := strconv.Atoi(id)

	if err != nil {
		return err
	}

	p, err := GetPlayerById(p1Id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, p)
}

func GetPlayerById(id int) (model.Player, error) {
	db := conn.Pool()
	defer db.Close()

	player := model.Player{}
	err := db.QueryRow(context.Background(), "SELECT * FROM player WHERE id=$1", id).Scan(
		&player.Id,
		&player.AccountId,
		&player.Play,
		&player.Hand,
		&player.Kitty,
		&player.Score,
		&player.Art,
	)

	if err != nil {
		return model.Player{}, err
	}

	return player, nil

}
