package container

import (
	"encoding/json"
	"queries"

	"cli/services"
	"cli/styles"
	"cli/utils"
	"cli/view/board"
	"cli/view/card"
	cliVO "cli/vo"
	"vo"

	tea "github.com/charmbracelet/bubbletea"
)

type Controller struct {
	cliVO.Controller
	subview cliVO.IController
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *Controller) Init() {
	ctrl.ChangeTab(0)
}

func (ctrl *Controller) Render() string {
	containerHeader := ctrl.View.Render()
	viewRender := ctrl.subview.Render()

	return containerHeader + "\n" + styles.WindowStyle.Render(viewRender)
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	subView := ctrl.subview

	subView.Update(msg)

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		resp := ctrl.ParseInput(msg)

		if resp == nil {
			break
		}

		cmds = append(cmds, func() tea.Msg {
			return resp
		})
	case vo.ChangeTabMsg:
		ctrl.ChangeTab(msg.TabIndex)
	}

	return tea.Batch(cmds...)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	containerModel := ctrl.Model.(*Model)
	containerView := ctrl.View.(*View)
	subView := ctrl.subview

	subView.ParseInput(msg)

	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "tab":
		containerView.ActiveTab = containerView.ActiveTab + 1
		if containerView.ActiveTab >= len(containerView.Tabs) {
			containerView.ActiveTab = 0
		}
		containerModel.State = containerView.Tabs[containerView.ActiveTab].TabState
		return vo.ChangeTabMsg{
			TabIndex: containerView.ActiveTab,
		}

	case "shift+tab":
		containerView.ActiveTab = containerView.ActiveTab - 1

		if containerView.ActiveTab < 0 {
			containerView.ActiveTab = len(containerView.Tabs) - 1
		}

		containerModel.State = containerView.Tabs[containerView.ActiveTab].TabState
		return vo.ChangeTabMsg{
			TabIndex: containerView.ActiveTab,
		}
	}

	return nil
}

func (ctrl *Controller) ChangeTab(tabIndex int) {
	containerModel := ctrl.Model.(*Model)

	gameDeck := ctrl.getGameDeck(containerModel)

	handModel := &cliVO.HandModel{
		CurrentTurnPlayerId: containerModel.Match.Players[0].ID,
		CardsToDisplay:      containerModel.Match.Players[0].Hand,
		SelectedCardIds:     []int32{},
		Deck:                gameDeck,
	}

	switch tabIndex {
	case 0:
		ctrl.subview = &board.Controller{
			Controller: cliVO.Controller{
				Model: board.Model{
					ViewModel: cliVO.ViewModel{
						Name: "Game",
					},
					LocalPlayerId: containerModel.Match.Players[0].ID,
					GameMatch:     containerModel.Match,
				},
			},
		}
	case 1:
		services.GetGampleCardsForMatch(containerModel.Match.ID)
		ctrl.subview = ctrl.CreateController("Play",
			handModel,
			gameDeck)
	case 2:
		services.GetGampleCardsForMatch(containerModel.Match.ID)
		ctrl.subview = ctrl.CreateController("Hand",
			handModel,
			gameDeck)
	case 3:
		services.GetGampleCardsForMatch(containerModel.Match.ID)
		ctrl.subview = ctrl.CreateController("Kitty",
			handModel,
			gameDeck)
	}

	ctrl.subview.Init()
}

func (ctrl *Controller) CreateController(name string, handModel *cliVO.HandModel, gameDeck *vo.GameDeck) cliVO.IController {
	m := &card.Model{
		ViewModel: &cliVO.ViewModel{
			Name: name,
		},
		ActiveSlotIndex:        0,
		HighlighedId:           0,
		HighlightedSlotIndexes: []int32{},
		Deck:                   gameDeck,
		HandModel:              handModel,
		SelectedCardId:         0,
	}

	v := &card.View{
		SelectedCardId: m.SelectedCardId,
		HandModel:      m.HandModel,
	}

	return &card.Controller{
		Controller: &cliVO.Controller{
			Model: m,
			View:  v,
		},
	}
}
func (ctrl *Controller) getGameDeck(containerModel *Model) *vo.GameDeck {
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

	return &vo.GameDeck{
		Deck:  deck,
		Cards: gameCards,
	}
}
