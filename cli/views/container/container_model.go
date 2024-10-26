package container

import (
	"cli/views"
	"model"
)

type ContainerModel struct {
	views.IViewModel
	Tabs   []views.Tab
	State  model.ViewState
	States []model.ViewState
	Match  *model.GameMatch
}
