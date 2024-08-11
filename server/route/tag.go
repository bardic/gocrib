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
// @Summary      Create new tag
// @Description
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param details body model.Tag true "New tag to save"
// @Success      200  {object}  model.Tag
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /tag/ [post]
func NewTag(c echo.Context) error {
	details := new(model.Tag)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Println(details)

	args := pgx.NamedArgs{
		"name": details.Name,
	}

	query := "INSERT INTO tags(name) VALUES (@name)"

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
// @Summary      Update tag by id
// @Description
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param details body model.Tag true "Item Object to save"
// @Success      200  {object}  model.Tag
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /tag/ [put]
func UpdateTag(c echo.Context) error {
	details := new(model.Tag)
	if err := c.Bind(details); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	args := pgx.NamedArgs{
		"id":   details.Id,
		"name": details.Name,
	}

	query := "UPDATE tags SET name = @name where id = @id"

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
// @Tags         tags
// @Accept       json
// @Produce      json
// @Param        tagid    query     string  true  "search for tag by id"'
// @Success      200  {object}  model.Tag
// @Failure      400  {object}  error
// @Failure      404  {object}  error
// @Failure      500  {object}  error
// @Router       /tag/ [get]
func GetTag(c echo.Context) error {
	t := c.Request().URL.Query().Get("tagid")

	db := model.Pool()
	defer db.Close()

	rows, err := db.Query(context.Background(), "SELECT * FROM tags WHERE id=$1", t)

	v := []model.Tag{}

	for rows.Next() {
		var item model.Tag
		err := rows.Scan(&item.Id, &item.Name)
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

	r, _ := json.Marshal(v[0])

	return c.JSON(http.StatusOK, string(r))
}
