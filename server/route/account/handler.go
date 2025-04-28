package account

import "github.com/bardic/gocrib/server/route"

type AccountHandler struct {
	route.BaseHandler
}

func NewHandler() *AccountHandler {
	return &AccountHandler{}
}
