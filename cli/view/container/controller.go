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
}

func (ctrl *Controller) GetState() cliVO.ControllerState {
	return cliVO.LoginControllerState
}

func (ctrl *Controller) Init() {
	ctrl.ChangeTab(0)
}

func (ctrl *Controller) Render() string {
	containerModel := ctrl.Model.(*Model)
	containerHeader := ctrl.View.Render()
	viewRender := containerModel.subview.Render()

	return containerHeader + "\n" + styles.WindowStyle.Render(viewRender)
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	containerModel := ctrl.Model.(*Model)
	subView := containerModel.subview

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
	subView := containerModel.subview

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

	handModel := ctrl.getHandModelForCardIds(
		containerModel.Match.Deckid,
		containerModel.Match.ID,
		containerModel.LocalPlayer.ID,
		containerModel.Match.Players[0].Hand,
	)
	switch tabIndex {
	case 0:
		containerModel.subview = &board.Controller{
			Controller: cliVO.Controller{
				Model: board.Model{
					ViewModel: cliVO.ViewModel{
						Name: "Game",
					},
					GameMatch:     containerModel.Match,
					LocalPlayerId: containerModel.LocalPlayer.ID,
				},
			},
		}
	case 1:

		containerModel.subview = ctrl.CreateController(
			"Play",
			handModel,

			containerModel.Match,
		)
	case 2:
		containerModel.subview = ctrl.CreateController(
			"Hand",
			handModel,
			containerModel.Match,
		)
	case 3:
		containerModel.subview = ctrl.CreateController(
			"Kitty",
			handModel,
			containerModel.Match,
		)
	}

	containerModel.subview.Init()
}

func (ctrl *Controller) getHandModelForCardIds(deckId, matchId, localPlayerId int32, cardIds []int32) *cliVO.HandVO {
	gameDeck := ctrl.getGameDeck(deckId, matchId)

	handModel := &cliVO.HandVO{
		LocalPlayerID: localPlayerId,
		CardIds:       cardIds,
		Deck:          gameDeck,
	}

	return handModel
}

func (ctrl *Controller) CreateController(name string, handModel *cliVO.HandVO, gameMatch *vo.GameMatch) cliVO.IController {
	m := &card.Model{
		ViewModel: &cliVO.ViewModel{
			Name: name,
		},
		ActiveSlotIndex:        0,
		HighlighedId:           0,
		HighlightedSlotIndexes: []int32{},
		Deck:                   handModel.Deck,
		HandVO:                 handModel,
		SelectedCardId:         0,
	}

	v := &card.View{
		ActiveCardId: m.SelectedCardId,
		HandVO:       m.HandVO,
	}

	return &card.Controller{
		Controller: &cliVO.Controller{
			Model: m,
			View:  v,
		},
		GameMatch: ctrl.Model.(*Model).Match,
	}
}
func (ctrl *Controller) getGameDeck(deckId, matchId int32) *vo.GameDeck {
	var deck *queries.Deck
	resp := services.GetDeckById(deckId)
	err := json.Unmarshal(resp.([]byte), &deck)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	var gameCards []queries.GetGameCardsForMatchRow
	resp = services.GetGampleCardsForMatch(matchId)
	err = json.Unmarshal(resp.([]byte), &gameCards)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	return &vo.GameDeck{
		Deck:  deck,
		Cards: gameCards,
	}
}
