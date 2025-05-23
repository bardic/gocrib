package container

import (
	"encoding/json"
	"time"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/view/board"
	"github.com/bardic/gocrib/cli/view/card"
	cliVO "github.com/bardic/gocrib/cli/vo"
	logger "github.com/bardic/gocrib/log"
	"github.com/bardic/gocrib/vo"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	view         *View
	model        *Model
	timer        timer.Model
	timerStarted bool
	tabIndex     int
	tabs         map[int]cliVO.IGameController
}

func NewController(accountID, matchID int) *Controller {
	l := logger.Get()
	defer l.Sync()

	var match *vo.Match

	resp := services.GetMatchByID(matchID)
	err := json.Unmarshal(resp.([]byte), &match)
	if err != nil {
		l.Sugar().Error(err)
	}

	player := &vo.Player{}

	var gameDeck vo.Deck
	resp = services.GetPlayerByForMatchAndAccount(matchID, accountID)
	err = json.Unmarshal(resp.([]byte), player)
	if err != nil {
		l.Sugar().Error(err)
	}

	resp = services.GetDeckByPlayIDAndMatchID(player.ID, matchID)
	err = json.Unmarshal(resp.([]byte), &gameDeck)
	if err != nil {
		l.Sugar().Error(err)
	}

	tabs := createTabs(match, player)
	ctrl := &Controller{
		model: NewModel(match, player, &gameDeck),
		view:  NewView(0, tabs),
	}

	ctrl.tabs = tabs

	ctrl.ChangeTab(0)
	return ctrl
}

func (ctrl *Controller) Init() {
}

func (ctrl *Controller) GetName() string {
	return "Container"
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *Controller) Render() string {
	viewRender := ctrl.model.Subcontroller.Render(ctrl.model.Gamematch, ctrl.model.GameDeck)

	return ctrl.view.Render() + "\n" + styles.WindowStyle.Render(viewRender)
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	l := logger.Get()
	defer l.Sync()

	var cmds []tea.Cmd

	if !ctrl.timerStarted {
		ctrl.timer = timer.NewWithInterval(time.Hour, time.Second*1)
		cmd := ctrl.timer.Init()
		cmds = append(cmds, cmd)
		ctrl.timerStarted = true
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		resp := ctrl.ParseInput(msg)

		if resp == nil {
			break
		}

		switch r := resp.(type) {
		case int:
			cmds = append(cmds, func() tea.Msg {
				return r
			})
		}
	case timer.TickMsg: // Polling update
		var cmd tea.Cmd
		var gameMatch *vo.Match
		var gameDeck *vo.Deck

		ctrl.timer, cmd = ctrl.timer.Update(msg)

		resp := services.GetMatchByID(ctrl.model.GetMatch().ID)
		err := json.Unmarshal(resp.([]byte), &gameMatch)
		if err != nil {
			l.Sugar().Error(err)
		}

		resp = services.GetDeckByPlayIDAndMatchID(ctrl.model.GetPlayer().ID, ctrl.model.GetMatch().ID)

		err = json.Unmarshal(resp.([]byte), &gameDeck)
		if err != nil {
			l.Sugar().Error(err)
		}

		ctrl.model.Gamematch = gameMatch
		ctrl.model.GameDeck = gameDeck
		cmds = append(cmds, cmd)
	case int:
		ctrl.ChangeTab(msg)
	}

	cmds = append(cmds, ctrl.model.Subcontroller.Update(msg))

	return tea.Batch(cmds...)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "tab":
		ctrl.view.ActiveTab++
		if ctrl.view.ActiveTab >= len(ctrl.view.Tabs) {
			ctrl.view.ActiveTab = 0
		}
		return ctrl.view.ActiveTab

	case "shift+tab":
		ctrl.view.ActiveTab--

		if ctrl.view.ActiveTab < 0 {
			ctrl.view.ActiveTab = len(ctrl.view.Tabs) - 1
		}

		return ctrl.view.ActiveTab
	default:
		return ctrl.model.Subcontroller.ParseInput(msg)
	}
}

func (ctrl *Controller) ChangeTab(idx int) {
	if ctrl.tabs == nil {
		ctrl.tabs = map[int]cliVO.IGameController{}
	}

	ctrl.tabIndex = idx

	val, ok := ctrl.tabs[idx]
	if ok {
		ctrl.model.Subcontroller = val
		// val.Update(msg)
		return
	}
}

func createTabs(gameMatch *vo.Match, player *vo.Player) map[int]cliVO.IGameController {
	return map[int]cliVO.IGameController{
		0: board.NewBoard(
			gameMatch,
			player,
		),
		1: card.NewController(
			"Play",
			gameMatch,
			player,
		),
		2: card.NewController(
			"Hand",
			gameMatch,
			player,
		),
		3: card.NewController(
			"Kitty",
			gameMatch,
			player,
		),
	}
}
