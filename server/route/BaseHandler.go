package route

import (
	"github.com/bardic/gocrib/server/store"
)

type BaseHandler struct {
	AccountStore store.AccountStore
	PlayerStore  store.PlayerStore
	DeckStore    store.DeckStore
	CardStore    store.CardStore
	MatchStore   store.MatchStore
}
