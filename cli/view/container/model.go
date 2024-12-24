package container

import (
	"github.com/bardic/gocrib/queries/queries"

	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"
)

type Model struct {
	cliVO.IModel
	Tabs        []cliVO.Tab
	State       vo.ViewState
	States      []vo.ViewState
	Match       *vo.GameMatch
	Subview     cliVO.IController
	LocalPlayer *queries.Player
	ActiveTab   int
}
