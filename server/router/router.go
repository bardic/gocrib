package router

import (
	"net/http"

	_ "github.com/bardic/gocrib/server/docs"
	"github.com/bardic/gocrib/server/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	// Add middleware to remove trailing slashes
	e.Pre(middleware.RemoveTrailingSlash())

	// Add swaggger router
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Add index router to return the version
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Cribbage  v0.0.1")
	})

	// Register v1 routes
	v2Routes(e.Group("/v2"))

	// Start Server
	e.Logger.Fatal(e.Start(":1323"))

	return e
}

func v2Routes(g *echo.Group) {
	h := route.NewHandler()

	a := g.Group("/account")
	a.POST("/:accountId", h.Login)

	m := g.Group("/match")
	m.GET("/open", h.GetOpenMatches)
	m.POST("/:accountId", h.NewMatch)
	m.GET("/:matchId/state", h.GetMatchState)
	m.GET("/:matchId", h.GetMatch)
	m.PUT("/:matchId/join/:accountId", h.JoinMatch)
	m.PUT("/:matchId/play", h.Play)
	m.PUT("/:matchId/cut/:deckIndex", h.CutDeck)
	m.GET("/:matchId/deck", h.GetDeck)
}
