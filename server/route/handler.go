package route

import (
	"fmt"
	"net/http"

	logger "github.com/bardic/gocrib/log"
	"github.com/bardic/gocrib/server/store"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	AccountStore store.AccountStore
	PlayerStore  store.PlayerStore
	DeckStore    store.DeckStore
	CardStore    store.CardStore
	MatchStore   store.MatchStore
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) BadParams(c echo.Context, details string, err error) error {
	e := addDetailsToError(details, err)
	return c.JSON(http.StatusBadRequest, e)
}

func (h *Handler) InternalError(c echo.Context, details string, err error) error {
	e := addDetailsToError(details, err)
	l := logger.Get()
	defer l.Sync()

	l.Sugar().Debugf("internal error %v", e)

	return c.JSON(http.StatusInternalServerError, e)
}

func (h *Handler) Ok(c echo.Context, msg any) error {
	return c.JSON(http.StatusOK, msg)
}

func addDetailsToError(details string, err error) error {
	return fmt.Errorf("%s : err = %w", details, err)
}
