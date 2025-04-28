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

	// db := conn.Pool()
	// defer db.Close()
	// q := queries.New(db)

	// b := BaseHandler{
	// 	AccountStore: store.AccountStore(&q),
	// 	PlayerStore:  *store.NewPlayerStore(&queries, ctx),
	// 	DeckStore:    *store.NewDeckStore(&queries, ctx),
	// 	CardStore:    *store.NewCardStore(&queries, ctx),
	// 	MatchStore:   *store.NewMatchStore(&queries, ctx),
	// }
	// handler := route.NewHandler()

	matchHandler := match.NewHandler()
	deckHandler := deck.NewHandler()
	playerHandler := player.NewHandler()
	accountHandler := account.NewHandler()

	g.GET("/open", matchHandler.GetOpenMatches)
	g.GET("/match/:matchId", matchHandler.GetMatch)
	g.POST("/match/:accountId", matchHandler.NewMatch)

	//Match
	matchGroup := g.Group("/match/:matchId")
	matchGroup.GET("/cards", matchHandler.GetMatchCardsForMatchId)
	matchGroup.GET("/account/:accountId", matchHandler.GetPlayerIdForMatchAndAccount)
	matchGroup.PUT("/cut/:cutIndex", matchHandler.CutDeck)
	matchGroup.PUT("/join/:accountId", matchHandler.JoinMatch)
	matchGroup.PUT("/pass", matchHandler.Pass)
	matchGroup.PUT("/currentPlayer/:playerId", matchHandler.UpdateCurrentPLayer)

	//Player
	matchGroup.GET("/player/:playerId", playerHandler.GetPlayer)
	playerGroup := g.Group("/player/:playerId")
	playerGroup.PUT("/to/:toPlayerId/kitty", playerHandler.UpdateKitty)
	playerGroup.PUT("/to/:toPlayerId/play", playerHandler.UpdatePlay)
	playerGroup.PUT("/ready", playerHandler.PlayerReady)
	playerGroup.GET("/deck", deckHandler.GetDeckByPlayerIdAndMatchId)

	//Deck
	matchGroup.GET("/deck", deckHandler.GetDeckByMatchId)
	deckGroup := matchGroup.Group("/deck")
	deckGroup.GET("/kitty", deckHandler.GetKitty)
	deckGroup.PUT("/shuffle", deckHandler.PutShuffle)

	//Account
	accountGroup := g.Group("/account")
	accountGroup.POST("/login/:accountId", accountHandler.Login)
}
