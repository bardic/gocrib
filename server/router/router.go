package router

import (
	"net/http"

	_ "github.com/bardic/gocrib/server/docs"

	"github.com/bardic/gocrib/server/route/account"
	"github.com/bardic/gocrib/server/route/deck"
	"github.com/bardic/gocrib/server/route/match"
	"github.com/bardic/gocrib/server/route/player"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Cribbage  v0.0.1")
	})

	v1Routes(e.Group("/v1"))

	e.Logger.Fatal(e.Start(":1323"))

	return e
}

func v1Routes(g *echo.Group) {
	matchHandler := match.NewHandler()
	deckHandler := deck.NewHandler()
	playerHandler := player.NewHandler()
	accountHandler := account.NewHandler()

	g.GET("/open", matchHandler.GetOpenMatches)
	g.GET("/match/:matchId", matchHandler.GetMatch)
	g.POST("/match/:accountId", matchHandler.NewMatch)

	// Match
	matchGroup := g.Group("/match/:matchId")
	matchGroup.GET("/cards", matchHandler.GetMatchCardsForMatchId)
	matchGroup.GET("/account/:accountId", matchHandler.GetPlayerIdForMatchAndAccount)
	matchGroup.PUT("/cut/:cutIndex", matchHandler.CutDeck)
	matchGroup.PUT("/join/:accountId", matchHandler.JoinMatch)

	// Player
	playerGroup := g.Group("/player/:playerId")
	playerGroup.PUT("/to/:toPlayerId/kitty", playerHandler.UpdateKitty)
	playerGroup.PUT("/to/:toPlayerId/play", playerHandler.UpdatePlay)

	// Deck
	matchGroup.GET("/deck", deckHandler.GetDeckByMatchId)

	// Account
	accountGroup := g.Group("/account")
	accountGroup.POST("/login/:accountId", accountHandler.Login)
}
