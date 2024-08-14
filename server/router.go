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

	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "judohippo" && password == "meow" {
			return true, nil
		}
		return false, nil
	}))

	//Card
	adminGroup.POST("/v1/card", route.NewCard)
	adminGroup.PUT("/v1/card", route.UpdateCard)
	adminGroup.DELETE("/v1/card", route.DeleteCard)
	//Player
	adminGroup.DELETE("/v1/player", route.DeletePlayer)
	//Match
	adminGroup.DELETE("/v1/match", route.DeleteMatch)
	//Lobby
	adminGroup.DELETE("/v1/lobby", route.DeleteLobby)
	//History
	adminGroup.DELETE("/v1/history", route.DeleteHistory)
	//Chat
	adminGroup.DELETE("/v1/chat", route.DeleteChat)

	playerGroup := e.Group("/player")

	//Card
	playerGroup.GET("/v1/card", route.GetCard)
	//Player
	playerGroup.GET("/v1/player", route.GetPlayer)
	playerGroup.POST("/v1/player", route.NewPlayer)
	playerGroup.PUT("/v1/player", route.UpdatePlayer)
	//Match
	playerGroup.GET("/v1/match", route.GetMatch)
	playerGroup.POST("/v1/match", route.NewMatch)
	playerGroup.PUT("/v1/match", route.UpdateMatch)
	// Lobby
	playerGroup.GET("/v1/lobby", route.GetLobby)
	playerGroup.POST("/v1/lobby", route.NewLobby)
	playerGroup.PUT("/v1/lobby", route.UpdateLobby)
	// History
	playerGroup.GET("/v1/history", route.GetHistory)
	playerGroup.POST("/v1/history", route.NewHistory)
	playerGroup.PUT("/v1/history", route.UpdateHistory)
	// Chat
	playerGroup.GET("/v1/chat", route.GetChat)
	playerGroup.POST("/v1/chat", route.NewChat)
	playerGroup.PUT("/v1/chat", route.UpdateChat)

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
