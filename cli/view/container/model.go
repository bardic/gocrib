package container

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.IModel
	Tabs          []cliVO.Tab
	State         vo.ViewState
	States        []vo.ViewState
	Match         *vo.GameMatch
	Subcontroller cliVO.IController
	LocalPlayer   vo.GamePlayer
	ActiveTab     int
}
