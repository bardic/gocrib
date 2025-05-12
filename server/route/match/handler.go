package match

import "github.com/bardic/gocrib/server/route"

type Handler struct {
	route.BaseHandler
}

func NewHandler() *Handler {
	return &Handler{}
}
