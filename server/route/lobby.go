package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bardic/cribbage/server/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new lobby
// @Description
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param details body model.Lobby true "lobby Object to save"
// @Success      200  {object}  model.Lobby
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /lobby/ [post]
func NewLobby(c echo.Context) error {
	details := new(model.Lobby)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseLobby(*details)

	query := "INSERT INTO lobbys(players, privateMatch, eloRangeMin, eloRangeMax, creatationDate) VALUES (@players, @privateMatch, @eloRangeMin, @eloRangeMax, @creatationDate)"

	db := model.Pool()
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
// @Summary      Update lobby by barcode
// @Description
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param details body model.Lobby true "lobby Object to save"
// @Success      200  {object}  model.Lobby
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /lobby/ [put]
func UpdateLobby(c echo.Context) error {
	details := new(model.Lobby)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseLobby(*details)

	query := "UPDATE lobbys SET players=@players, privateMatch=@privateMatch, eloRangeMin=@eloRangeMin, eloRangeMax=@eloRangeMax, creatationDate=@creatationDate id = @id"

	db := model.Pool()
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
// @Summary      Get lobby by barcode
// @Description
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for lobby by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.Lobby
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /lobby/ [get]
func GetLobby(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM lobbys WHERE id=$1", id)

	v := []model.Lobby{}

	for rows.Next() {
		var lobby model.Lobby

		err := rows.Scan(&lobby.Id, &lobby.Players, &lobby.PrivateMatch, &lobby.EloRangeMin, &lobby.EloRangeMax, &lobby.CreatationDate)
		if err != nil {
			fmt.Println(err)
			return err
		}

		v = append(v, lobby)
	}

	if err != nil {
		return err
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get lobby by barcode
// @Description
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for lobby by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.Lobby
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /lobby/ [delete]
func DeleteLobby(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// lobby, _ := getlobby(b, s)

	return c.JSON(http.StatusOK, nil)
}

func parseLobby(details model.Lobby) pgx.NamedArgs {
	return pgx.NamedArgs{
		"Accounts":       details.Players,
		"PrivateMatch":   details.PrivateMatch,
		"EloRangeMin":    details.EloRangeMin,
		"EloRangeMax":    details.EloRangeMax,
		"CreatationDate": details.CreatationDate,
	}
}
