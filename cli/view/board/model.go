package board

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.ViewModel
	AccountId *int
	CutIndex  string
	GameMatch *vo.GameMatch
}
