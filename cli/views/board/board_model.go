package board

import (
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
)

type BoardModel struct {
	views.ViewModel
	Account       *queries.Account
	CutIndex      string
	GameMatch     *model.GameMatch
	LocalPlayerId int32
}
