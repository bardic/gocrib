package container

import (
	"encoding/json"
	"queries"

	"cli/services"
	"cli/styles"
	"cli/utils"
	"cli/views"
	"cli/views/board"
	"cli/views/card"
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
	gc.View = &ContainerView{
		ActiveTab: 0,
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

	return containerHeader + "\n" + styles.WindowStyle.Render(viewRender)
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
		containerView.ActiveTab = containerView.ActiveTab + 1
		containerModel.State = containerView.tabs[containerView.ActiveTab].TabState

	case "shift+tab":
		containerView.ActiveTab = containerView.ActiveTab - 1
		containerModel.State = containerView.tabs[containerView.ActiveTab].TabState
	}
	return model.ChangeTabMsg{
		TabIndex: containerView.ActiveTab,
	}
}

func (gc *ContainerController) ChangeTab(tabIndex int) {
	containerModel := gc.Model.(*ContainerModel)

	gameDeck := gc.getGameDeck(containerModel)

	handModel := &views.HandModel{
		CurrentTurnPlayerId: containerModel.Match.Players[0].ID,
		CardsToDisplay:      containerModel.Match.Players[0].Hand,
		SelectedCardIds:     []int32{},
		Deck:                gameDeck,
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
		gc.CreateController("Play",
			containerModel.Match.Players[0].Kitty,
			handModel,
			gameDeck)
	case 2:
		services.GetGampleCardsForMatch(containerModel.Match.ID)
		gc.CreateController("Hand",
			containerModel.Match.Players[0].Hand,
			handModel,
			gameDeck)
	case 3:
		services.GetGampleCardsForMatch(containerModel.Match.ID)
		gc.CreateController("Kitty",
			containerModel.Match.Players[0].Kitty,
			handModel,
			gameDeck)
	}

	gc.subview.Init()
}

func (gc *ContainerController) CreateController(name string, cards []int32, handModel *views.HandModel, gameDeck *model.GameDeck) views.IController {
	return &card.CardController{
		Controller: &views.Controller{
			Model: &card.CardModel{
				ViewModel: &views.ViewModel{
					Name: name,
				},
				ActiveSlotIdx:       0,
				HighlighedId:        0,
				HighlightedSlotIdxs: []int{},
				Cards:               cards,
				Deck:                gameDeck,
			},
			View: &card.CardView{
				SelectedCardId: 0,
				HandModel:      handModel,
			},
		},
	}
}
func (gc *ContainerController) getGameDeck(containerModel *ContainerModel) *model.GameDeck {
	var deck *queries.Deck
	resp := services.GetDeckById(containerModel.Match.Deckid)
	err := json.Unmarshal(resp.([]byte), &deck)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	var gameCards []queries.GetGameCardsForMatchRow
	resp = services.GetGampleCardsForMatch(containerModel.Match.ID)
	err = json.Unmarshal(resp.([]byte), &gameCards)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	return &model.GameDeck{
		Deck:  deck,
		Cards: gameCards,
	}
}
