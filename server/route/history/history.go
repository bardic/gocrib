package route

import (
	"context"
	"fmt"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/cribbagev2/model"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Create new history
// @Description
// @Tags         historys
// @Accept       json
// @Produce      json
// @Param details body model.History true "history Object to save"
// @Success      200  {object}  model.History
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /history/ [post]
func NewHistory(c echo.Context) error {
	details := new(model.History)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseHistory(*details)
	query := "INSERT INTO historys(matchId, matchCompletetionDate, winners, losers) VALUES (@matchId, @matchCompletetionDate, @currentOwner, @originalOwner, @winners, @losers)"

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
// @Summary      Update history by barcode
// @Description
// @Tags         historys
// @Accept       json
// @Produce      json
// @Param details body model.History true "history Object to save"
// @Success      200  {object}  model.History
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /history/ [put]
func UpdateHistory(c echo.Context) error {
	details := new(model.History)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := parseHistory(*details)

	query := "UPDATE historys SET matchId = @matchId, matchCompletetionDate = @matchCompletetionDate, winners = @winners, losers=@losers, storename=@storeName, storeNeighborhood=@storeNeighborhood, tags=@tags where barcode = @barcode AND storeId = @storeId"

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
// @Summary      Get history by barcode
// @Description
// @Tags         historys
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for history by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.History
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /history/ [get]
func GetHistory(c echo.Context) error {
	matchId := c.Request().URL.Query().Get("matchId")

	db := conn.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM historys WHERE matchId=$1", matchId)

	v := []model.History{}

	for rows.Next() {
		var history model.History

		err := rows.Scan(&history.MatchId, &history.MatchCompletetionDate, &history.Winners, &history.Losers)
		if err != nil {
			fmt.Println(err)
			return &echo.BindingError{}
		}

		v = append(v, history)
	}

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, v)
}

// Create godoc
// @Summary      Get history by barcode
// @Description
// @Tags         historys
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for history by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.History
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /history/ [delete]
func DeleteHistory(c echo.Context) error {
	// b := c.Request().URL.Query().Get("barcode")
	// s := c.Request().URL.Query().Get("storeId")

	// history, _ := gethistory(b, s)

	return c.JSON(http.StatusOK, nil)
}

func parseHistory(details model.History) pgx.NamedArgs {
	return pgx.NamedArgs{
		"matchId":               details.MatchId,
		"matchCompletetionDate": details.MatchCompletetionDate,
		"winners":               details.Winners,
		"losers":                details.Losers,
	}
}
