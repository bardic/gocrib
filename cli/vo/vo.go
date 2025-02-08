package vo

import (
	"github.com/bardic/gocrib/vo"

	tea "github.com/charmbracelet/bubbletea"
)

type ControllerState int

const (
	LoginControllerState ControllerState = iota
	LobbyControllerState
	BoardControllerState
)

type IController interface {
	Render(gameMatch *vo.GameMatch) string
	ParseInput(tea.KeyMsg) tea.Msg
	Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd
	GetState() ControllerState
	GetModel() IModel
}

type IView interface {
	Render(hand []int) string
	Init()
	BuildHeader() string
	BuildFooter() string
}

type IHandView interface {
	IView
	BuildHeader() string
	BuildFooter() string
}

type HandVO struct {
	LocalPlayerID int
	CardIds       []int
	Deck          *vo.GameDeck
}

type IModel interface {
	GetMatch() *vo.GameMatch
	GetPlayer() *vo.GamePlayer
}

type ViewModel struct {
	Name string
}

func (m ViewModel) GetMatch() *vo.GameMatch {
	return m.GetMatch()
}

func (m ViewModel) GetPlayer() *vo.GamePlayer {
	return m.GetPlayer()
}

type GameController struct {
	View       IView
	Model      IModel
	Controller IController
}

func (ctrl *GameController) GetModel() IModel {
	return ctrl.Model
}

func (ctrl *GameController) GetView() IView {
	return ctrl.View
}

func (ctrl *GameController) SetModel(model IModel) {
	ctrl.Model = model
}

func (ctrl *GameController) SetView(view IView) {
	ctrl.View = view
}

type Tab struct {
	Title    string
	TabState vo.ViewState
}
