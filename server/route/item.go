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
// @Summary      Create new item
// @Description
// @Tags         items
// @Accept       json
// @Produce      json
// @Param details body model.Item true "Item Object to save"
// @Success      200  {object}  model.Item
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /item/ [post]
func NewItem(c echo.Context) error {
	details := new(model.Item)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Println(details)

	args := pgx.NamedArgs{
		"name":              details.Name,
		"cost":              details.Cost,
		"weight":            details.Weight,
		"unit":              details.Unit,
		"barcode":           details.Barcode,
		"storeId":           details.StoreId,
		"storeName":         details.StoreName,
		"storeNeighborhood": details.StoreNeighborhood,
		"tags":              details.Tags,
	}
	fmt.Println("What the hell is happening 1")
	fmt.Println(args)

	query := "INSERT INTO items(name, cost, weight, unit, barcode, storeid, storename, storeneighborhood, tags) VALUES (@name, @cost, @weight, @unit, @barcode, @storeId, @storeName, @storeNeighborhood, @tags)"

	fmt.Println("What the hell is happening")
	fmt.Println(query)

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
// @Summary      Update item by barcode
// @Description
// @Tags         items
// @Accept       json
// @Produce      json
// @Param details body model.Item true "Item Object to save"
// @Success      200  {object}  model.Item
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /item/ [put]
func UpdateItem(c echo.Context) error {
	details := new(model.Item)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := pgx.NamedArgs{
		"name":              details.Name,
		"cost":              details.Cost,
		"weight":            details.Weight,
		"unit":              details.Unit,
		"barcode":           details.Barcode,
		"storeId":           details.StoreId,
		"storeName":         details.StoreName,
		"storeNeighborhood": details.StoreNeighborhood,
		"tags":              details.Tags,
	}

	query := "UPDATE items SET name = @name, cost = @cost, weight = @weight, unit=@unit, storename=@storeName, storeNeighborhood=@storeNeighborhood, tags=@tags where barcode = @barcode AND storeId = @storeId"

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
// @Summary      Get item by barcode
// @Description
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "search for item by barcode"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.Item
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /item/ [get]
func GetItem(c echo.Context) error {
	b := c.Request().URL.Query().Get("barcode")
	s := c.Request().URL.Query().Get("storeId")

	item, _ := getItem(b, s)

	return c.JSON(http.StatusOK, item)
}

func getItem(barcode string, storeid string) (string, error) {
	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM items WHERE barcode=$1 and storeId=$2", barcode, storeid)

	v := []model.Item{}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.Id, &item.Name, &item.Cost, &item.Weight, &item.Unit, &item.Barcode, &item.StoreId, &item.StoreName, &item.StoreNeighborhood, &item.Tags, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			return "", err
		}

		v = append(v, item)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	r, _ := json.Marshal(v)

	return string(r), nil
}
