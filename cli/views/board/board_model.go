package board

import (
	"cli/views"
	"model"
	"queries"
)

type BoardModel struct {
	views.ViewModel
	Account       *queries.Account
	CutIndex      string
	GameMatch     *model.GameMatch
	LocalPlayerId int32
}
