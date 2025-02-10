package lobby

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.ViewModel
}

func (m Model) GetMatch() *vo.GameMatch {
	return m.Gamematch
}
