package player

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/gocrib/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new player
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body int true "player Object to save"
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [post]
func NewPlayer(c echo.Context) error {
	var details int
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	p, err := NewPlayerQuery(details)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, p)

}

func NewPlayerQuery(accountId int) (*model.Player, error) {

	args := pgx.NamedArgs{
		"accountId": accountId,
		"hand":      []int{},
		"kitty":     []int{},
		"play":      []int{},
		"score":     0,
		"art":       "default.png",
	}

	query := `INSERT INTO player (
			accountid,
			hand,
			play,
			kitty,
			score,
			art
		) VALUES (
			@accountId,
			@hand,
			@play,
			@kitty,
			@score,
			@art
		)
		RETURNING id`

	db := conn.Pool()
	defer db.Close()

	var playerId int
	err := db.QueryRow(
		context.Background(),
		query,
		args).Scan(&playerId)

	if err != nil {
		return nil, err
	}

	p := model.Player{}
	p.Id = playerId
	p.AccountId = accountId

	return &p, nil
}
