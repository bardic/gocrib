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

	matchGroup := g.Group("/match")
	matchGroup.GET("/player", player.GetPlayer)
	matchGroup.POST("/player", player.NewPlayer)
	matchGroup.PUT("/player/kitty", player.UpdateKitty)
	matchGroup.PUT("/player/ready", player.PlayerReady)
	matchGroup.PUT("/player/play", match.UpdatePlay)
	matchGroup.GET("", match.GetMatch)
	matchGroup.GET("/state", match.GetMatchState)
	matchGroup.GET("/open", match.GetOpenMatches)
	matchGroup.POST("", match.NewMatch)
	matchGroup.PUT("/join", match.JoinMatch)
	matchGroup.PUT("/cut", match.CutDeck)
	matchGroup.GET("/cards", match.GetMatchCardsForMatchId)

	g.GET("/deck", deck.GetDeckByMatchId)

	//Account
	accountGroup := g.Group("/account")
	accountGroup.POST("/login", account.Login)
}
