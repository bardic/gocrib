package route

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bardic/cribbage/server/model"

	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Get items from store
// @Description
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        storeId    query     string  true  "Store which to get items from"'
// @Success      200  {object}  model.Items
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /items/fromStore [get]
func GetItemsFromStore(c echo.Context) error {
	s := c.Request().URL.Query().Get("storeId")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM items WHERE storeId=$1", s)

	v := []model.Item{}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.Id, &item.Name, &item.Cost, &item.Weight, &item.Unit, &item.Barcode, &item.StoreId, &item.StoreName, &item.StoreNeighborhood, &item.Tags, &item.CreatedAt, &item.UpdatedAt)
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

	return c.JSON(http.StatusOK, v)
}

// Create godoc
// @Summary      Get items related to barcode
// @Description
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        barcode    query     string  true  "barcode to get related items to"'
// @Param        storeId    query     string  true  "Store in which the barcode was found"'
// @Success      200  {object}  model.Item
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /items/related  [get]
func GetRelatedItemsForBarcode(c echo.Context) error {
	b := c.Request().URL.Query().Get("barcode")
	s := c.Request().URL.Query().Get("storeId")

	var data []model.Item

	r, _ := getItem(b, s)

	err := json.Unmarshal([]byte(r), &data)

	if err != nil {
		fmt.Println(err)
		return err
	}

	db := model.Pool()
	defer db.Close()

	var sb strings.Builder

	sb.WriteString("'(")
	for _, v := range data {
		for i, tag := range v.Tags {
			sb.WriteString(fmt.Sprintf("%%%s%%", strconv.Itoa(tag)))
			if i < len(v.Tags)-1 {
				sb.WriteString("|")
			}
		}
	}
	sb.WriteString(")'")

	q := fmt.Sprintf("SELECT * FROM items WHERE CAST(tags as text) SIMILAR TO %s", sb.String())
	rows, err := db.Query(context.Background(), q)

	if err != nil {
		fmt.Println(err)
	}

	v := []model.Item{}

	for rows.Next() {
		var item model.Item
		err := rows.Scan(&item.Id, &item.Name, &item.Cost, &item.Weight, &item.Unit, &item.Barcode, &item.StoreId, &item.StoreName, &item.StoreNeighborhood, &item.Tags, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return err
		}

		v = append(v, item)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	g, _ := json.Marshal(v)
	return c.JSON(http.StatusOK, string(g))

}
