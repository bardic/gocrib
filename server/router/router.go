// Description: Main router for the server, sets up the echo server and routes
package router

import (
	_ "github.com/bardic/gocrib/server/docs"

	"github.com/bardic/gocrib/server/route/account"
	"github.com/bardic/gocrib/server/route/deck"
	"github.com/bardic/gocrib/server/route/match"
	"github.com/bardic/gocrib/server/route/player"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Router struct
type Router struct {
}

// Creates new Router for grocitcdn boycot, configs middleware and API paths
func (r *Router) New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())

	/******
	*
	* Admin
	*
	********/
	v1 := e.Group("/v1")
	v1Routes(v1)

	/******
	*
	* SWAGGER
	*
	********/

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Cribbage  v0.0.1")
	})

	e.Logger.Fatal(e.Start(":1323"))

	return e
}

func v1Routes(g *echo.Group) {
	adminGroup := g.Group("/admin")
	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "judohippo" && password == "meow" {
			return true, nil
		}
		return false, nil
	}))

	g.GET("/open", match.GetOpenMatches)
	g.GET("/match/:matchId", match.GetMatch)
	// g.GET("/deck/", deck.GetDeckByMatchId)

	g.POST("/match/:accountId", match.NewMatch)

	//Match
	matchGroup := g.Group("/match/:matchId")
	matchGroup.GET("/cards", match.GetMatchCardsForMatchId)
	matchGroup.PUT("/cut/:cutIndex", match.CutDeck)
	matchGroup.PUT("/join/:accountId", match.JoinMatch)
	matchGroup.PUT("/deal", match.Deal)
	matchGroup.PUT("/determinefirst", match.DetermineFirst)
	matchGroup.PUT("/pass", match.Pass)
	matchGroup.GET("/deck", deck.GetDeckByMatchId)
	matchGroup.PUT("/currentPlayer/:playerId", match.UpdateCurrentPLayer)
	// matchGroup.GET("/kitty", match.GetKitty)

	//Player
	matchGroup.GET("/player/:playerId", player.GetPlayer)
	playerGroup := matchGroup.Group("/player/:playerId")
	playerGroup.PUT("/kitty", player.UpdateKitty)

	playerGroup.PUT("/play", player.UpdatePlay)
	playerGroup.PUT("/ready", player.PlayerReady)

	//Deck
	deckGroup := matchGroup.Group("/deck")
	deckGroup.GET("/kitty", deck.GetKitty)
	deckGroup.PUT("/shuffle", deck.PutShuffle)

	//Account
	accountGroup := g.Group("/account")
	accountGroup.POST("/login", account.Login)
}
