package container

import (
	cliVO "cli/vo"
	"queries"
	"vo"
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
