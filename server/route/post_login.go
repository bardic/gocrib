package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Login route
// @Description Login route for account - takes an account id and returns the account details
// @Tags		account
// @Accept		json
// @Produce	json
// @Param		accountId		path		int	true	"account id"'
// @Success	200		{object}	vo.Account
// @Failure	400		{object}	error
// @Failure	500		{object}	error
// @Router		/account/{accountId} [post]
func (h *Handler) Login(c echo.Context) error {
	accountID, err := strconv.Atoi(c.Param("accountId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("can't convert param accountId: %w", err))
	}

	account, err := h.AccountStore.GetAccountByID(c, accountID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to get account by id: %w", err))
	}

	return c.JSON(http.StatusOK, account)
}
