package gameContainer

import (
	"encoding/json"
	"time"

	"github.com/bardic/gocrib/vo"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/container"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*container.Controller
	timer        timer.Model
	timerStarted bool
}

func (ctrl *Controller) Update(msg tea.Msg, oldGameMatch *vo.GameMatch) tea.Cmd {
	var cmds []tea.Cmd

	if !ctrl.timerStarted {
		ctrl.timer = timer.NewWithInterval(time.Hour, time.Second*1)
		cmd := ctrl.timer.Init()
		cmds = append(cmds, cmd)
		ctrl.timerStarted = true
	}

	switch msg := msg.(type) {
	case timer.TickMsg: // Polling update
		var cmd tea.Cmd
		var gameMatch *vo.GameMatch
		ctrl.timer, cmd = ctrl.timer.Update(msg)

		resp := services.GetMatchById(ctrl.Controller.Model.(*container.Model).Match.ID)
		err := json.Unmarshal(resp.([]byte), &gameMatch)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		ctrl.Controller.Model.(*container.Model).Match = gameMatch

		cmds = append(cmds, cmd)

	}

	cmd := ctrl.Controller.Update(msg, ctrl.Controller.Model.(*container.Model).Match)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
