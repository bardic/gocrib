package lobby

import (
	"github.com/bardic/gocrib/cli/views"
)

type LobbyModel struct {
	views.ViewModel
	ActiveMatchId int32
	AccountId     int32
}
