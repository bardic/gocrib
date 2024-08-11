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
// @Summary      Create new store position
// @Description
// @Tags         storePosition
// @Accept       json
// @Produce      json
// @Param details body model.StorePosition true "StorePosition to save"
// @Success      200  {object}  model.StorePosition
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /store/ [post]
func NewStorePosition(c echo.Context) error {
	details := new(model.StorePosition)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := pgx.NamedArgs{
		"storeId": details.StoreId,
		"lat":     details.Lat,
		"long":    details.Long,
	}

	fmt.Println(args)

	query := "INSERT INTO storeposition(storeId, lat, long) VALUES (@storeId, @lat, @long)"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return c.JSON(http.StatusOK, "meow")
}

// Create godoc
// @Summary      Update storePosition by storeID
// @Description
// @Tags         storePosition
// @Accept       json
// @Produce      json
// @Param details body model.StorePosition true "Update storePosition for storeId"
// @Success      200  {object}  model.StorePosition
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /store/ [put]
func UpdateStorePosition(c echo.Context) error {
	details := new(model.StorePosition)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := pgx.NamedArgs{
		"storeId": details.StoreId,
		"lat":     details.Lat,
		"long":    details.Long,
	}

	query := "UPDATE storeposition SET lat = @lat, long = @long where storeid = @storeId"

	db := model.Pool()
	defer db.Close()

	_, err := db.Exec(
		context.Background(),
		query,
		args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return c.JSON(http.StatusOK, "meow")
}

// Create godoc
// @Summary      Get store by storeId
// @Description
// @Tags         storePosition
// @Accept       json
// @Produce      json
// @Param        storeId    query     string  true  "StorePosition for storeID"'
// @Success      200  {object}  model.StorePosition
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /store/ [get]
func GetStorePositionByStoreId(c echo.Context) error {
	s := c.Request().URL.Query().Get("storeId")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM storeposition WHERE storeid = $1", s)

	v := []model.StorePosition{}

	for rows.Next() {
		var item model.StorePosition
		err := rows.Scan(&item.Id, &item.StoreId, &item.Lat, &item.Long, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			return err
		}

		v = append(v, item)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	r, _ := json.Marshal(v)
	return c.JSON(http.StatusOK, string(r))
}

// Create godoc
// @Summary      Get store by storeId
// @Description
// @Tags         storePosition
// @Accept       json
// @Produce      json
// @Param        lat    query     int  true  "StorePosition for storeID"'
// @Param        long    query     int  true  "StorePosition for storeID"'
// @Success      200  {object}  model.StorePosition
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /store/byPosition [get]
func GetStorePositionByPosition(c echo.Context) error {
	lat := c.Request().URL.Query().Get("lat")
	long := c.Request().URL.Query().Get("long")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT *, point($1, $2) <@>  (point(long, lat)::point) as distance FROM storeposition ORDER BY distance;", long, lat)

	v := []model.StorePosition{}

	for rows.Next() {
		var item model.StorePosition
		err := rows.Scan(&item.Id, &item.StoreId, &item.Lat, &item.Long, &item.CreatedAt, &item.UpdatedAt, &item.Distance)
		if err != nil {
			fmt.Println(err)
			return err
		}

		v = append(v, item)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	r, _ := json.Marshal(v)
	return c.JSON(http.StatusOK, string(r))
}
