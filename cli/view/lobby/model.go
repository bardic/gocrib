package lobby

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
)

type Model struct {
	cliVO.ViewModel
	ActiveMatchId int32
	AccountId     int32
}
