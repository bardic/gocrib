// Desc: Login route for account
package account

import (
	"context"
	"net/http"
	"time"

	"github.com/bardic/gocrib/queries/queries"

	conn "github.com/bardic/gocrib/server/db"

	"github.com/labstack/echo/v4"
)

// Login user via id
//
//	@Summary	Login
//	@Description
//	@Tags		account
//	@Accept		json
//	@Produce	json
//	@Param		details	body		int	true	"id to login with"
//	@Success	200		{object}	queries.Account
//	@Failure	400		{object}	error
//	@Failure	500		{object}	error
//	@Router		/account/login/ [post]
func Login(c echo.Context) error {
	id := new(int32)
	if err := c.Bind(id); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db := conn.Pool()
	defer db.Close()
	q := queries.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a, err := q.GetAccount(ctx, *id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, a)
}
