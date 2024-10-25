package container

import (
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
)

type ContainerModel struct {
	views.IViewModel
	ActiveTab int
	Tabs      []views.Tab
	State     model.ViewState
	States    []model.ViewState
	Match     *model.GameMatch
}
