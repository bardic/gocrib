package lobby

import (
	"cli/views"
)

type LobbyModel struct {
	views.ViewModel
	ActiveMatchId int32
	AccountId     int32
}
