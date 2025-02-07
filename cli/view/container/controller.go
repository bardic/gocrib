package container

import (
	"encoding/json"

	"github.com/bardic/gocrib/queries/queries"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/cli/view/board"
	"github.com/bardic/gocrib/cli/view/card"
	cliVO "github.com/bardic/gocrib/cli/vo"
	"github.com/bardic/gocrib/vo"

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

func (ctrl *Controller) Render(gamematch *vo.GameMatch) string {
	containerModel := ctrl.Model.(*Model)

	cardIds := []int{}

	for _, card := range containerModel.LocalPlayer.Hand {
		cardIds = append(cardIds, *card.Cardid)
	}

	containerHeader := ctrl.View.Render(cardIds)
	viewRender := containerModel.Subcontroller.Render(containerModel.Match)

	return containerHeader + "\n" + styles.WindowStyle.Render(viewRender)
}

func (ctrl *Controller) Update(msg tea.Msg, gameMatch *vo.GameMatch) tea.Cmd {
	// var cmd tea.Cmd
	var cmds []tea.Cmd
	containerModel := ctrl.Model.(*Model)
	subView := containerModel.Subcontroller

	if gameMatch != nil {
		cmd := subView.Update(msg, gameMatch)
		cmds = append(cmds, cmd)
	}

	// cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg: //User input
		resp := ctrl.ParseInput(msg)

		if resp == nil {
			break
		}

		switch r := resp.(type) {
		case vo.ChangeTabMsg:
			cmds = append(cmds, func() tea.Msg {
				return r
			})
		}

	case vo.ChangeTabMsg:
		ctrl.ChangeTab(msg.TabIndex)
	}

	return tea.Batch(cmds...)
}

func (ctrl *Controller) ParseInput(msg tea.KeyMsg) tea.Msg {
	containerModel := ctrl.Model.(*Model)
	containerView := ctrl.View.(*View)
	//subView := containerModel.Subview

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
	default:
		containerModel.Subcontroller.ParseInput(msg)
	}

	return msg
}

func (ctrl *Controller) ChangeTab(tabIndex int) {
	containerModel := ctrl.Model.(*Model)
	// deckId := containerModel.Match.Deckid

	// hand := getPlayerHand(containerModel.LocalPlayer.ID, containerModel.Match.Players)

	switch tabIndex {
	case 0:
		containerModel.Subcontroller = &board.Controller{
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
		containerModel.Subcontroller = ctrl.CreateController(
			"Play",
			containerModel.Match.Gamestate,
			ctrl.getHandModelForCardIds(
				*containerModel.LocalPlayer.ID,
				*containerModel.Match.ID,
				utils.IdFromCards(containerModel.LocalPlayer.Play),
			),
		)
	case 2:
		containerModel.Subcontroller = ctrl.CreateController(
			"Hand",
			containerModel.Match.Gamestate,
			ctrl.getHandModelForCardIds(
				*containerModel.LocalPlayer.ID,
				*containerModel.Match.ID,
				utils.IdFromCards(containerModel.LocalPlayer.Hand), //THIS IS THE WRONG PLAYER~!!!!!
			),
		)
	case 3:
		containerModel.Subcontroller = ctrl.CreateController(
			"Kitty",
			containerModel.Match.Gamestate,
			ctrl.getHandModelForCardIds(
				*containerModel.LocalPlayer.ID,
				*containerModel.Match.ID,
				utils.IdFromCards(containerModel.LocalPlayer.Kitty),
			),
		)
	}

	containerModel.Subcontroller.Init()

}

func getPlayerHand(playerId *int, players []vo.GamePlayer) []int {
	for _, p := range players {
		if p.ID == playerId {
			return utils.IdFromCards(p.Hand)
		}
	}

	return []int{}
}

func (ctrl *Controller) getHandModelForCardIds(localPlayerId, matchId int, cardIds []int) *cliVO.HandVO {
	//gameDeck := ctrl.getGameDeckForMatchId(matchId)
	gameDeck := ctrl.getDeckByPlayerIdAndMatchId(localPlayerId, matchId)

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
		SelectedCardIds: []int{},
		Deck:            handModel.Deck,
		HandVO:          handModel,
		State:           currentState,
	}

	v := &card.View{
		// ActiveCardId:   m.ActiveSlotIndex,
		Deck: m.HandVO.Deck,
		// ActivePlayerId: m.HandVO.LocalPlayerID,
		// MatchId:        ctrl.Model.(*Model).Match.ID,
		// GameState:      ctrl.Model.(*Model).Match.Gamestate,
	}

	return &card.Controller{
		Controller: &cliVO.Controller{
			Model: m,
			View:  v,
		},
		GameMatch: ctrl.Model.(*Model).Match,
	}
}

// func (ctrl *Controller) getGameDeckForMatchId(matchId int) *vo.GameDeck {
// 	var deck *vo.GameDeck
// 	resp := services.GetDeckByMatchId(matchId)
// 	err := json.Unmarshal(resp.([]byte), &deck)
// 	if err != nil {
// 		utils.Logger.Sugar().Error(err)
// 	}

// 	return deck
// }

func (ctrl *Controller) getDeckByPlayerIdAndMatchId(playerId, matchId int) *vo.GameDeck {
	var deck *vo.GameDeck

	resp := services.GetDeckByPlayIdAndMatchId(playerId, matchId)
	err := json.Unmarshal(resp.([]byte), &deck)
	if err != nil {
		utils.Logger.Sugar().Error(err)
	}

	return deck
}
