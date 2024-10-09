package account

import (
	"context"
	"net/http"

	conn "github.com/bardic/cribbage/server/db"
	"github.com/bardic/gocrib/model"
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
	id := new(int)
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()

	var name string
	err := db.QueryRow(context.Background(), "SELECT name FROM accounts WHERE id=$1", id).Scan(&name)

	if err != nil {
		return err
	}

	a := model.Account{
		Id:   *id,
		Name: name,
	}

	return c.JSON(http.StatusOK, a)
}
