package game

import (
	cliVO "github.com/bardic/gocrib/cli/vo"
)

type Controller struct {
	View       cliVO.IView
	Model      cliVO.IModel
	Controller cliVO.IController
}

func (ctrl *Controller) GetModel() cliVO.IModel {
	return ctrl.Model
}

func (ctrl *Controller) GetView() cliVO.IView {
	return ctrl.View
}

func (ctrl *Controller) SetModel(model cliVO.IModel) {
	ctrl.Model = model
}

func (ctrl *Controller) SetView(view cliVO.IView) {
	ctrl.View = view
}
