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
// @Router       /player/lobby/ [post]
func NewLobby(c echo.Context) error {
	details := new(model.Lobby)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseLobby(*details)

	query := "INSERT INTO lobby(players, privateMatch, eloRangeMin, eloRangeMax) VALUES (@players, @privateMatch, @eloRangeMin, @eloRangeMax) RETURNING id"

	db := model.Pool()
	defer db.Close()

	row := db.QueryRow(
		context.Background(),
		query,
		args)

	var lobby model.Lobby
	err := row.Scan(&lobby.Id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, lobby.Id)
}

// Create godoc
// @Summary      Update lobby by id
// @Description
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param details body model.Lobby true "lobby Object to save"
// @Success      200  {object}  model.Lobby
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/lobby/ [put]
func UpdateLobby(c echo.Context) error {
	details := new(model.Lobby)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseLobby(*details)

	query := "UPDATE lobby SET players=@players, privateMatch=@privateMatch, eloRangeMin=@eloRangeMin, eloRangeMax=@eloRangeMax where id = @id"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "Lobby Updated")
}

// Create godoc
// @Summary      Get lobby by id
// @Description
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for id by barcode"'
// @Success      200  {object}  model.Lobby
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/lobby/ [get]
func GetLobby(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM lobby WHERE id=$1", id)

	if err != nil {
		return err
	}

	v := []model.Lobby{}

	for rows.Next() {
		var lobby model.Lobby

		err := rows.Scan(&lobby.Id, &lobby.Players, &lobby.CreatationDate, &lobby.PrivateMatch, &lobby.EloRangeMin, &lobby.EloRangeMax)
		if err != nil {
			fmt.Println(err)
			return err
		}

		v = append(v, lobby)
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      delete lobby by id
// @Description
// @Tags         lobbys
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "delete for lobby by id"'
// @Success      200  {object}  model.Lobby
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/lobby/ [delete]
func DeleteLobby(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	//DELETE FROM lobby WHERE where id=id;

	id := c.Request().URL.Query().Get("id")

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(context.Background(), "DELETE FROM lobby WHERE id=$1", id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func parseLobby(details model.Lobby) pgx.NamedArgs {
	return pgx.NamedArgs{
		"players":        details.Players,
		"privateMatch":   details.PrivateMatch,
		"eloRangeMin":    details.EloRangeMin,
		"eloRangeMax":    details.EloRangeMax,
		"creatationDate": details.CreatationDate,
	}
}
