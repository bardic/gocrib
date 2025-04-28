package match

import "github.com/bardic/gocrib/server/route"

type MatchHandler struct {
	route.BaseHandler
}

func NewHandler() *MatchHandler {
	return &MatchHandler{}
}
