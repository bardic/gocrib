package route

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestMatchFor(t *testing.T) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/player/match/?id=1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	fmt.Println(GetMatch(c))
}
