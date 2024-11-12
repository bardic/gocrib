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
	*cliVO.Controller
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
	viewRender := containerModel.Subview.Render()

	return containerHeader + "\n" + styles.WindowStyle.Render(viewRender)
}

func (ctrl *Controller) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	containerModel := ctrl.Model.(*Model)
	subView := containerModel.Subview

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
	subView := containerModel.Subview

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
	deckId := containerModel.Match.Deckid

	hand := getPlayerHand(containerModel.LocalPlayer.ID, containerModel.Match.Players)

	switch tabIndex {
	case 0:
		containerModel.Subview = &board.Controller{
			Controller: &cliVO.Controller{
				Model: &board.Model{
					ViewModel: cliVO.ViewModel{
						Name: "Game",
					},
					GameMatch:     containerModel.Match,
					LocalPlayerId: containerModel.LocalPlayer.ID,
				},
				View: &board.View{
					Match:         containerModel.Match,
					LocalPlayerId: containerModel.LocalPlayer.ID,
					State:         containerModel.Match.Gamestate,
				},
			},
		}
	case 1:
		containerModel.Subview = ctrl.CreateController(
			"Play",
			containerModel.Match.Gamestate,
			ctrl.getHandModelForCardIds(
				containerModel.LocalPlayer.ID,
				deckId,
				containerModel.LocalPlayer.Play,
			),
		)
	case 2:
		containerModel.Subview = ctrl.CreateController(
			"Hand",
			containerModel.Match.Gamestate,
			ctrl.getHandModelForCardIds(
				containerModel.LocalPlayer.ID,
				deckId,
				hand, //THIS IS THE WRONG PLAYER~!!!!!
			),
		)
	case 3:
		containerModel.Subview = ctrl.CreateController(
			"Kitty",
			containerModel.Match.Gamestate,
			ctrl.getHandModelForCardIds(
				containerModel.LocalPlayer.ID,
				deckId,
				containerModel.LocalPlayer.Kitty,
			),
		)
	}

	containerModel.Subview.Init()

}

func getPlayerHand(playerId int32, players []*queries.Player) []int32 {
	for _, p := range players {
		if p.ID == playerId {
			return p.Hand
		}
	}

	return []int32{}
}

func (ctrl *Controller) getHandModelForCardIds(localPlayerId, deckId int32, cardIds []int32) *cliVO.HandVO {
	gameDeck := ctrl.getGameDeck(deckId)

	handModel := &cliVO.HandVO{
		LocalPlayerID: localPlayerId,
		CardIds:       cardIds,
		Deck:          gameDeck,
	}

	return handModel
}

func (ctrl *Controller) CreateController(name string, currentState queries.Gamestate, handModel *cliVO.HandVO) cliVO.IController {
	m := &card.Model{
		ViewModel: &cliVO.ViewModel{
			Name: name,
		},
		ActiveSlotIndex: 0,
		SelectedCardIds: []int32{},
		Deck:            handModel.Deck,
		HandVO:          handModel,
		State:           currentState,
	}

	v := &card.View{
		ActiveCardId:   m.ActiveSlotIndex,
		HandVO:         m.HandVO,
		ActivePlayerId: m.HandVO.LocalPlayerID,
		MatchId:        ctrl.Model.(*Model).Match.ID,
		GameState:      ctrl.Model.(*Model).Match.Gamestate,
	}

	return &card.Controller{
		Controller: &cliVO.Controller{
			Model: m,
			View:  v,
		},
		GameMatch: ctrl.Model.(*Model).Match,
	}
}
func (ctrl *Controller) getGameDeck(deckId int32) *vo.GameDeck {
	var deck *vo.GameDeck
	resp := services.GetDeckById(deckId)
	err := json.Unmarshal(resp.([]byte), &deck)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	return deck
}
