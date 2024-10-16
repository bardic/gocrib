package main

import (
	_ "github.com/bardic/cribbage/server/docs"

	"github.com/bardic/cribbage/server/route/account"
	"github.com/bardic/cribbage/server/route/admin/card"
	"github.com/bardic/cribbage/server/route/deck"
	"github.com/bardic/cribbage/server/route/gameplaycard"
	"github.com/bardic/cribbage/server/route/match"
	"github.com/bardic/cribbage/server/route/player"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

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

	//Card
	adminGroup.POST("/card", card.NewCard)
	adminGroup.PUT("/card", card.UpdateCard)
	adminGroup.DELETE("/card", card.DeleteCard)
	//Player
	adminGroup.DELETE("/player", player.DeletePlayer)
	//Match
	adminGroup.DELETE("/match", match.DeleteMatch)

	// //History
	// adminGroup.DELETE("/history", route.DeleteHistory)
	// //Chat
	// adminGroup.DELETE("/chat", route.DeleteChat)

	playerGroup := g.Group("/player")

	//Card
	// playerGroup.GET("/card", card.GetCard)
	// playerGroup.GET("/allcards", route.GetAllCards)
	//Player
	playerGroup.GET("/gameplaycards", gameplaycard.GetGameplayCards)
	playerGroup.GET("/player", player.GetPlayer)
	playerGroup.POST("/player", player.NewPlayer)
	playerGroup.PUT("/player", player.UpdatePlayer)
	playerGroup.PUT("/kitty", player.UpdateKitty)
	playerGroup.PUT("/ready", player.PlayerReady)
	// playerGroup.PUT("/play", player.UpdatePlay)
	//Match
	playerGroup.GET("/match", match.GetMatch)
	playerGroup.GET("/match/state", match.GetMatchState)
	// playerGroup.GET("/matches", match.GetMatches)
	playerGroup.GET("/matches/open", match.GetOpenMatches)
	playerGroup.POST("/match", match.NewMatch)
	playerGroup.PUT("/match", match.UpdateMatch)
	playerGroup.PUT("/match/join", match.JoinMatch)
	playerGroup.PUT("/match/cut", match.CutDeck)

	playerGroup.GET("/match/deck", deck.GetDeck)
	// History
	// playerGroup.GET("/history", route.GetHistory)
	// playerGroup.POST("/history", route.NewHistory)
	// playerGroup.PUT("/history", route.UpdateHistory)
	// // Chat
	// playerGroup.GET("/chat", route.GetChat)
	// playerGroup.POST("/chat", route.NewChat)
	// playerGroup.PUT("/chat", route.UpdateChat)

	//Account
	accountGroup := g.Group("/account")
	accountGroup.POST("/login", account.Login)

	// gameGroup := g.Group("/game")
	// gameGroup.POST("/playCard", player.PlayCard)
}
