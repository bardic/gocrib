package board

import (
	cliVO "cli/vo"
	"queries"
	"vo"
)

type Model struct {
	cliVO.ViewModel
	Account       *queries.Account
	CutIndex      string
	GameMatch     *vo.GameMatch
	LocalPlayerId int32
}
