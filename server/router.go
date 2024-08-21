package main

import (
	_ "github.com/bardic/cribbage/server/docs"
	"github.com/bardic/cribbage/server/route"
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
	adminGroup.POST("/card", route.NewCard)
	adminGroup.PUT("/card", route.UpdateCard)
	adminGroup.DELETE("/card", route.DeleteCard)
	//Player
	adminGroup.DELETE("/player", route.DeletePlayer)
	//Match
	adminGroup.DELETE("/match", route.DeleteMatch)

	//History
	adminGroup.DELETE("/history", route.DeleteHistory)
	//Chat
	adminGroup.DELETE("/chat", route.DeleteChat)

	playerGroup := g.Group("/player")

	//Card
	playerGroup.GET("/card", route.GetCard)
	//Player
	playerGroup.GET("/player", route.GetPlayer)
	playerGroup.POST("/player", route.NewPlayer)
	playerGroup.PUT("/player", route.UpdatePlayer)
	//Match
	playerGroup.GET("/match", route.GetMatch)
	playerGroup.POST("/match", route.NewMatch)
	playerGroup.PUT("/match", route.UpdateMatch)
	// Lobby
	playerGroup.GET("/lobby", route.GetLobby)
	playerGroup.POST("/lobby", route.NewLobby)
	playerGroup.PUT("/lobby", route.UpdateLobby)
	playerGroup.DELETE("/lobby", route.DeleteLobby)
	// History
	playerGroup.GET("/history", route.GetHistory)
	playerGroup.POST("/history", route.NewHistory)
	playerGroup.PUT("/history", route.UpdateHistory)
	// Chat
	playerGroup.GET("/chat", route.GetChat)
	playerGroup.POST("/chat", route.NewChat)
	playerGroup.PUT("/chat", route.UpdateChat)

	gameGroup := g.Group("/game")
	gameGroup.POST("/playCard", route.PlayCard)
}
