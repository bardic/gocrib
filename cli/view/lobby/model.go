package lobby

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.ViewModel
	ActiveMatchId *int
	AccountId     *int
}

func (m Model) GetMatch() *vo.GameMatch {
	return m.GetMatch()
}
