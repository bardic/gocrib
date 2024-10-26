package lobby

import (
	cliVO "cli/vo"
)

type Model struct {
	cliVO.ViewModel
	ActiveMatchId int32
	AccountId     int32
}
