package deck

import "github.com/bardic/gocrib/server/route"

type DeckHandler struct {
	route.BaseHandler
}

func NewHandler() *DeckHandler {
	return &DeckHandler{}
}
