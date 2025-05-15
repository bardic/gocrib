package router

import (
	"net/http"

	_ "github.com/bardic/gocrib/server/docs"

	"github.com/bardic/gocrib/server/route/account"
	"github.com/bardic/gocrib/server/route/match"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
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

// func v1Routes(g *echo.Group) {
// 	matchHandler := match.NewHandler()
// 	deckHandler := deck.NewHandler()
// 	playerHandler := player.NewHandler()
// 	accountHandler := account.NewHandler()
//
// 	g.GET("/open", matchHandler.GetOpenMatches)
// 	g.GET("/match/:matchId", matchHandler.GetMatch)
// 	g.POST("/match/:accountId", matchHandler.NewMatch)
//
// 	// Match
// 	matchGroup := g.Group("/match/:matchId")
// 	matchGroup.GET("/cards", matchHandler.GetMatchCardsForMatchID)
// 	matchGroup.GET("/account/:accountId", matchHandler.GetPlayerIDForMatchAndAccount)
// 	matchGroup.PUT("/cut/:cutIndex", matchHandler.CutDeck)
// 	matchGroup.PUT("/join/:accountId", matchHandler.JoinMatch)
//
// 	// Player
// 	playerGroup := g.Group("/player/:playerId")
// 	playerGroup.PUT("/to/:toPlayerId/kitty", playerHandler.UpdateKitty)
// 	playerGroup.PUT("/to/:toPlayerId/play", playerHandler.UpdatePlay)
//
// 	// Deck
// 	matchGroup.GET("/deck", deckHandler.GetDeckByMatchID)
//
// 	// Account
// 	accountGroup := g.Group("/account")
// 	accountGroup.POST("/login/:accountId", accountHandler.Login)
// }

func v2Routes(g *echo.Group) {
	matchHandler := match.NewHandler()
	// deckHandler := deck.NewHandler()
	// playerHandler := player.NewHandler()
	accountHandler := account.NewHandler()

	a := g.Group("/account")
	a.POST("/:accountId", accountHandler.Login)

	m := g.Group("/match")
	m.GET("/open", matchHandler.GetOpenMatches)
	m.POST("/:accountId", matchHandler.NewMatch)
	m.GET("/:matchId/state", matchHandler.GetMatchState)
	m.GET("/:matchId", matchHandler.GetMatch)
	m.PUT("/:matchId/join/:accountId", matchHandler.JoinMatch)
	m.PUT("/:matchId/play", matchHandler.Play)
	m.PUT("/:matchId/cut/:deckIndex", matchHandler.CutDeck)
	m.GET("/:matchId/deck", matchHandler.GetDeck)
}
