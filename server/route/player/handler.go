package player

import "github.com/bardic/gocrib/server/route"

type PlayerHandler struct {
	route.BaseHandler
}

func NewHandler() *PlayerHandler {
	return &PlayerHandler{}
}
