package gameContainer

import (
	"encoding/json"
	"time"
	"vo"

	"cli/services"
	"cli/utils"
	"cli/view/container"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	*container.Controller
	timer        timer.Model
	timerStarted bool
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmd = ctrl.Controller.Update(msg)
	cmds = append(cmds, cmd)

	if !ctrl.timerStarted {
		ctrl.timer = timer.NewWithInterval(time.Hour, time.Second*1)
		cmd := ctrl.timer.Init()
		cmds = append(cmds, cmd)
		ctrl.timerStarted = true
	}

	switch msg := msg.(type) {
	case timer.TickMsg: // Polling update
		var cmd tea.Cmd
		ctrl.timer, cmd = ctrl.timer.Update(msg)
		var gameMatch vo.GameMatch
		resp := services.GetPlayerMatch(ctrl.Controller.Model.(*container.Model).Match.ID)
		err := json.Unmarshal(resp.([]byte), &gameMatch)

		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		ctrl.Controller.Model.(*container.Model).Match = &gameMatch

		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}
