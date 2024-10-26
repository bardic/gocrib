package container

import (
	cliVO "cli/vo"
	"vo"
)

type Model struct {
	cliVO.IViewModel
	Tabs      []cliVO.Tab
	State     vo.ViewState
	States    []vo.ViewState
	Match     *vo.GameMatch
	ActiveTab int
}
