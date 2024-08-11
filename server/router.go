package main

import (
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
	* ITEM
	*
	********/

	e.POST("/v1/item", route.NewItem)
	e.PUT("/v1/item", route.UpdateItem)
	e.GET("/v1/item", route.GetItem)

	/******
	*
	* TAG
	*
	********/

	e.POST("/v1/tag", route.NewTag)
	e.PUT("/v1/tag", route.UpdateTag)
	e.GET("/v1/tag", route.GetTag)

	/******
	*
	* ITEM
	*
	********/

	e.GET("/v1/items/fromStore", route.GetItemsFromStore)
	e.GET("/v1/items/related", route.GetRelatedItemsForBarcode)

	/******
	*
	* STORE
	*
	********/

	e.POST("/v1/store", route.NewStorePosition)
	e.PUT("/v1/store", route.UpdateStorePosition)
	e.GET("/v1/store", route.GetStorePositionByStoreId)
	e.GET("/v1/store/byPosition", route.GetStorePositionByPosition)

	/******
	*
	* SWAGGER
	*
	********/

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "grocit v0.0.2")
	})

	e.Logger.Fatal(e.Start(":1323"))

	return e
}
