package match

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"queries"
	"vo"

	"github.com/labstack/echo/v4"
)

func TestNewMatch(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/player/match", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("details", vo.MatchRequirements{
		EloRangeMin: 0,
		EloRangeMax: 500,
		IsPrivate:   false,
	})

	err := NewMatch(c)

	if err != nil {
		t.Fatalf(`meow`)
	}

	var m queries.Match
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m.ID == 0 {
		t.Fatalf(`meow`)
	}
}

func TestMatch(t *testing.T) {
	// e := echo.New()
	// req := httptest.NewRequest(http.MethodGet, "/player/match/?id=1", nil)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// c := e.NewContext(req, rec)

	// err := GetMatch(c)

	// if err != nil {
	// 	t.Fatalf(`meow`)
	// }

	// req2 := httptest.NewRequest(http.MethodGet, "/player/match/?id=1", nil)
	// req2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec2 := httptest.NewRecorder()
	// c2 := e.NewContext(req, rec)

	// var match vo.Match
	// json.Unmarshal([]byte(rec2.Body.String()), &match)

	// err = GetDeck(c2)

	// if err != nil {
	// 	t.Fatalf(`meow`)
	// }
}

func TestDeck(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/player/match/?id=1", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetMatch(c)

	if err != nil {
		t.Fatalf(`meow`)
	}

	//err = GetDeck(c.)
}

func TestNoMatch(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/player/deck/?id=0", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetMatch(c)

	if err == nil {
		t.Fatalf(`Failed to return error when no match found`)
	}
}
