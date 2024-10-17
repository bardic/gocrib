package account

import (
	"context"
	"net/http"

	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	conn "github.com/bardic/gocrib/server/db"
	"github.com/labstack/echo/v4"
)

// Create godoc
// @Summary      Login
// @Description
// @Tags         account
// @Accept       json
// @Produce      json
// @Param details body int true "id to login with"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  error
// @Failure      500  {object}  error
// @Router       /account/login/ [post]
func Login(c echo.Context) error {
	id := new(int32)
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()

	ctx := context.Background()

	q := queries.New(db)

	name, err := q.GetAccount(ctx, *id)

	// var name string
	// err := db.QueryRow("SELECT name FROM accounts WHERE id=$1", id).Scan(&name)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	a := model.Account{
		Id:   int(*id),
		Name: name,
	}

	return c.JSON(http.StatusOK, a)
}
