package route

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMatch(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/player/match/?id=1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetMatch(c)

	if err != nil {
		t.Fatalf(`meow`)
	}
}

func TestNoMatch(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/player/match/?id=0", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetMatch(c)

	if err == nil {
		t.Fatalf(`Failed to return error when no match found`)
	}
}
