package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbagev2/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new player
// @Description
// @Tags         players
// @Accept       json
// @Produce      json
// @Param details body model.Player true "player Object to save"
// @Success      200  {object}  model.Player
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/player/ [post]
func NewPlayer(c echo.Context) error {
	details := new(model.Player)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parsePlayer(*details)

	query := "INSERT INTO player(hand, kitty, score, art) VALUES (@hand, @kitty, @score, @art)"

	db := conn.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

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

	err := UpdatePlayerQuery(*details)

	if err != nil {
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, nil)
}

func UpdatePlayerQuery(player model.Player) error {
	args := parsePlayer(player)

	query := "UPDATE player SET hand = @hand, kitty = @kitty, score = @score, art = @art where id = @id"

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

	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM player WHERE id=$1", id)

	v := []model.Player{}

	for rows.Next() {
		var player model.Player

		err := rows.Scan(&player.Id, &player.Hand, &player.Kitty, &player.Score, &player.Art)
		if err != nil {
			fmt.Println(err)
			return err
		}

		v = append(v, player)
	}

	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

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
// @Router       /player/player/ [delete]
func DeletePlayer(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// player, _ := getplayer(b, s)

	return c.JSON(http.StatusOK, nil)
}

func parsePlayer(details model.Player) pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":    details.Id,
		"hand":  details.Hand,
		"kitty": details.Kitty,
		"score": details.Score,
		"art":   details.Art,
	}
}
