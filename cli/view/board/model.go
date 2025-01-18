package board

import (
	"github.com/bardic/gocrib/queries/queries"

	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.ViewModel
	Account       *queries.Account
	CutIndex      string
	GameMatch     *vo.GameMatch
	LocalPlayerId int
}
