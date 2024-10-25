package container

import (
	"encoding/json"

	"cli/services"
	"cli/utils"
	"cli/views"
	"cli/views/board"
	"cli/views/kitty"
	"cli/views/play"
	"cli/views/playerhand"
	"model"

	tea "github.com/charmbracelet/bubbletea"
)

type ContainerController struct {
	views.Controller
	subview views.IController
}

func (gc *ContainerController) GetState() views.ControllerState {
	return views.LoginControllerState
}

func (gc *ContainerController) Init() {
	containerModel := gc.Model.(*ContainerModel)
	gc.View = &ContainerView{
		activeTab: containerModel.ActiveTab,
		tabs: []views.Tab{
			{
				Title:    "Board",
				TabState: model.BoardView,
			},
			{
				Title:    "Play",
				TabState: model.PlayView,
			},
			{
				Title:    "Hand",
				TabState: model.HandView,
			},
			{
				Title:    "Kitty",
				TabState: model.KittyView,
			},
		},
	}

	gc.ChangeTab(0)
}

func (gc *ContainerController) Render() string {
	containerHeader := gc.View.Render()
	viewRender := gc.subview.Render()

	return containerHeader + "\n" + viewRender
}

func (v *ContainerController) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		resp := v.ParseInput(msg)

		if resp == nil {
			break
		}

		cmds = append(cmds, func() tea.Msg {
			return resp
		})
	case model.ChangeTabMsg:
		v.ChangeTab(msg.TabIndex)
	}

	return tea.Batch(cmds...)
}

func (v *ContainerController) ParseInput(msg tea.KeyMsg) tea.Msg {
	containerModel := v.Model.(*ContainerModel)
	containerView := v.View.(*ContainerView)

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "tab":
		containerModel.ActiveTab = containerModel.ActiveTab + 1
		containerModel.State = containerView.tabs[containerModel.ActiveTab].TabState

	case "shift+tab":
		containerModel.ActiveTab = containerModel.ActiveTab - 1
		containerModel.State = containerView.tabs[containerModel.ActiveTab].TabState
	}
	return model.ChangeTabMsg{
		TabIndex: containerModel.ActiveTab,
	}
}

func (gc *ContainerController) ChangeTab(tabIndex int) {
	containerModel := gc.Model.(*ContainerModel)

	handModel := views.HandModel{
		CurrentTurnPlayerId: containerModel.Match.Players[0].ID,
		CardsToDisplay:      containerModel.Match.Players[0].Hand,
		SelectedCardIds:     []int32{},
		Deck:                &model.GameDeck{},
	}

	switch tabIndex {
	case 0:
		gc.subview = &board.BoardController{
			Controller: views.Controller{
				Model: board.BoardModel{
					ViewModel: views.ViewModel{
						Name: "Game",
					},
					LocalPlayerId: containerModel.Match.Players[0].ID,
					GameMatch:     containerModel.Match,
				},
			},
		}
	case 1:
		services.GetGampleCardsForMatch(containerModel.Match.ID)

		var deck *model.GameDeck
		resp := services.GetDeckById(containerModel.Match.Deckid)
		err := json.Unmarshal(resp.([]byte), &deck)
		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		gc.subview = &play.PlayController{
			Controller: &views.Controller{
				Model: &play.PlayModel{
					ViewModel: views.ViewModel{
						Name: "Play",
					},
					ActiveSlotIdx:       0,
					HighlighedId:        0,
					HighlightedSlotIdxs: []int{},
					Cards:               containerModel.Match.Players[0].Hand,
					Deck:                deck,
				},
				View: &play.PlayView{
					SelectedCardId: 0,
					HandModel:      handModel,
				},
			},
		}
	case 2:
		services.GetGampleCardsForMatch(containerModel.Match.ID)

		var deck *model.GameDeck
		resp := services.GetDeckById(containerModel.Match.Deckid)
		err := json.Unmarshal(resp.([]byte), &deck)
		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		gc.subview = &play.PlayController{
			Controller: &views.Controller{
				Model: &play.PlayModel{
					ViewModel: views.ViewModel{
						Name: "Play",
					},
					ActiveSlotIdx:       0,
					HighlighedId:        0,
					HighlightedSlotIdxs: []int{},
					Cards:               containerModel.Match.Players[0].Hand,
					Deck:                deck,
				},
				View: &playerhand.PlayerHandView{
					SelectedCardId: 0,
					HandModel:      handModel,
				},
			},
		}
	case 3:
		services.GetGampleCardsForMatch(containerModel.Match.ID)

		var deck *model.GameDeck
		resp := services.GetDeckById(containerModel.Match.Deckid)
		err := json.Unmarshal(resp.([]byte), &deck)
		if err != nil {
			utils.Logger.Sugar().Error(err)
		}

		gc.subview = &play.PlayController{
			Controller: &views.Controller{
				Model: &play.PlayModel{
					ViewModel: views.ViewModel{
						Name: "Play",
					},
					ActiveSlotIdx:       0,
					HighlighedId:        0,
					HighlightedSlotIdxs: []int{},
					Cards:               containerModel.Match.Players[0].Kitty,
					Deck:                deck,
				},
				View: &kitty.KittyView{
					SelectedCardId: 0,
					HandModel:      handModel,
				},
			},
		}
	}

	gc.subview.Init()
}
