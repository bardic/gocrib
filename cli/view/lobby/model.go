package lobby

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
)

type Model struct {
	cliVO.ViewModel
	ActiveMatchId *int
	AccountId     *int
}
