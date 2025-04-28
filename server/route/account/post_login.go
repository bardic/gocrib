// Desc: Login route for account
package account

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Login user via id
//
//	@Summary	Login
//	@Description
//	@Tags		account
//	@Accept		json
//	@Produce	json
//	@Param		accountId		path		int	true	"account id"'
//	@Success	200		{object}	queries.Account
//	@Failure	400		{object}	error
//	@Failure	500		{object}	error
//	@Router		/account/login/{accountId} [post]
func (h *AccountHandler) Login(c echo.Context) error {
	accountId, err := strconv.Atoi(c.Param("accountId"))

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	account, err := h.AccountStore.GetAccountById(c, &accountId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, account)
}
