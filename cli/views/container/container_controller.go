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

func (cc *ContainerController) GetState() views.ControllerState {
	return views.LoginControllerState
}

func (cc *ContainerController) Init() {
	cc.ChangeTab(0)
}

func (cc *ContainerController) Render() string {
	containerHeader := cc.View.Render()
	viewRender := cc.subview.Render()

	return containerHeader + "\n" + styles.WindowStyle.Render(viewRender)
}

func (cc *ContainerController) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		resp := cc.ParseInput(msg)

		if resp == nil {
			break
		}

		cmds = append(cmds, func() tea.Msg {
			return resp
		})
	case model.ChangeTabMsg:
		cc.ChangeTab(msg.TabIndex)
	}

	return tea.Batch(cmds...)
}

func (cc *ContainerController) ParseInput(msg tea.KeyMsg) tea.Msg {
	containerModel := cc.Model.(*ContainerModel)
	containerView := cc.View.(*ContainerView)

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "tab":
		containerView.ActiveTab = containerView.ActiveTab + 1
		containerModel.State = containerView.Tabs[containerView.ActiveTab].TabState

	case "shift+tab":
		containerView.ActiveTab = containerView.ActiveTab - 1
		containerModel.State = containerView.Tabs[containerView.ActiveTab].TabState
	}
	return model.ChangeTabMsg{
		TabIndex: containerView.ActiveTab,
	}
}

func (cc *ContainerController) ChangeTab(tabIndex int) {
	containerModel := cc.Model.(*ContainerModel)

	gameDeck := cc.getGameDeck(containerModel)

	handModel := &views.HandModel{
		CurrentTurnPlayerId: containerModel.Match.Players[0].ID,
		CardsToDisplay:      containerModel.Match.Players[0].Hand,
		SelectedCardIds:     []int32{},
		Deck:                gameDeck,
	}

	switch tabIndex {
	case 0:
		cc.subview = &board.BoardController{
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
		cc.subview = cc.CreateController("Play",
			handModel,
			gameDeck)
	case 2:
		services.GetGampleCardsForMatch(containerModel.Match.ID)
		cc.subview = cc.CreateController("Hand",
			handModel,
			gameDeck)
	case 3:
		services.GetGampleCardsForMatch(containerModel.Match.ID)
		cc.subview = cc.CreateController("Kitty",
			handModel,
			gameDeck)
	}

	cc.subview.Init()
}

func (cc *ContainerController) CreateController(name string, handModel *views.HandModel, gameDeck *model.GameDeck) views.IController {
	m := &card.CardModel{
		ViewModel: &views.ViewModel{
			Name: name,
		},
		ActiveSlotIdx:       0,
		HighlighedId:        0,
		HighlightedSlotIdxs: []int{},
		Deck:                gameDeck,
		HandModel:           handModel,
		SelectedCardId:      0,
	}

	v := &card.CardView{
		SelectedCardId: m.SelectedCardId,
		HandModel:      m.HandModel,
	}

	return &card.CardController{
		Controller: &views.Controller{
			Model: m,
			View:  v,
		},
	}
}
func (cc *ContainerController) getGameDeck(containerModel *ContainerModel) *model.GameDeck {
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
