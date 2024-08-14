package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bardic/cribbage/server/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new match
// @Description
// @Tags         matchs
// @Accept       json
// @Produce      json
// @Param details body model.match true "match Object to save"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /match/ [post]
func NewMatch(c echo.Context) error {
	details := new(model.Match)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseMatch(*details)

	query := "INSERT INTO matchs(value, suit, currentOwner, originalOwner, state, art) VALUES (@value, @suit, @currentOwner, @originalOwner, @state, @art)"

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
// @Summary      Update match by barcode
// @Description
// @Tags         matchs
// @Accept       json
// @Produce      json
// @Param details body model.match true "match Object to save"
// @Success      200  {object}  model.match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /match/ [put]
func UpdateMatch(c echo.Context) error {
	details := new(model.Match)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseMatch(*details)

	query := "UPDATE matchs SET name = @name, cost = @cost, weight = @weight, unit=@unit, storename=@storeName, storeNeighborhood=@storeNeighborhood, tags=@tags where barcode = @barcode AND storeId = @storeId"

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
// @Summary      Get match by barcode
// @Description
// @Tags         matchs
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for match by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /match/ [get]
func GetMatch(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM matchs WHERE id=$1", id)

	v := []model.Match{}

	for rows.Next() {
		var match model.Match

		err := rows.Scan(&match.Id, &match.LobbyId, &match.BoardId, &match.CurrentPlayerTurn, &match.TurnPassTimestamps)
		if err != nil {
			fmt.Println(err)
			return err
		}

		v = append(v, match)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	r, _ := json.Marshal(v)

	return c.JSON(http.StatusOK, r)
}

// Create godoc
// @Summary      Get match by barcode
// @Description
// @Tags         matchs
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for match by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /match/ [delete]
func DeleteMatch(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// match, _ := getmatch(b, s)

	return c.JSON(http.StatusOK, nil)
}

func parseMatch(details model.Match) pgx.NamedArgs {
	return pgx.NamedArgs{
		"lobbyId":            details.LobbyId,
		"boardId":            details.BoardId,
		"currentPlayerTurn":  details.CurrentPlayerTurn,
		"turnPassTimestamps": details.TurnPassTimestamps,
	}
}
