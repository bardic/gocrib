package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bardic/cribbage/server/model"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new match card
// @Description
// @Tags         match_card
// @Accept       json
// @Produce      json
// @Param details body model.Match true "match Object to save"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/card [post]
func NewMatchCard(c echo.Context) error {
	details := new(model.Match)
	fmt.Print(time.Now())
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseMatch(*details)

	query := "INSERT INTO match(lobbyId, currentPlayerTurn, art) VALUES (@lobbyId, @currentPlayerTurn, @art)"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "meow")
}

// Create godoc
// @Summary      Update match by barcode
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param details body model.Match true "match Object to save"
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/card [put]
func UpdateMatchCard(c echo.Context) error {
	details := new(model.Match)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseMatch(*details)

	if err := updateMatchQuery(args); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "meow")
}

// Create godoc
// @Summary      Get match by barcode
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /player/match/card [get]
func GetMatchCard(c echo.Context) error {
	id := c.Request().URL.Query().Get("id")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM match WHERE id=$1", id)

	v := []model.Match{}

	for rows.Next() {
		var match model.Match

		err := rows.Scan(&match.Id, &match.LobbyId, &match.CardsInPlay, &match.CutGameCardId, &match.CurrentPlayerTurn, &match.TurnPassTimestamps, &match.Art)
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

	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get match by barcode
// @Description
// @Tags         match
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "search for match by barcode"'
// @Success      200  {object}  model.Match
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /admin/match/card [delete]
func DeleteMatchCard(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// match, _ := getmatch(b, s)

	return c.JSON(http.StatusOK, nil)
}
