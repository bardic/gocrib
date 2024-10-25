package hand

import (
	"slices"

	"github.com/bardic/gocrib/cli/services"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
	"github.com/bardic/gocrib/queries"
	tea "github.com/charmbracelet/bubbletea"
)

type HandController struct {
	views.Controller
}

func (gc *HandController) GetState() views.ControllerState {
	return views.LoginControllerState
}

func (hc *HandController) Init() {

}
func (hc *HandController) Render() string {
	return ""
}

func (hc *HandController) InitView() {

}
func (hc *HandController) ParseInput(msg tea.KeyMsg) tea.Msg {
	handModel := hc.Model.(*HandModel)
	switch msg.String() {
	case "ctrl+c", "q":
		return tea.Quit()
	case "enter", "view_update":
		switch handModel.State {
		case queries.GamestateDiscardState:
			services.PutKitty(model.HandModifier{
				MatchId:  handModel.MatchId,
				PlayerId: handModel.PlayerId,
				CardIds:  handModel.Cards,
			})
		case queries.GamestatePlayState:
			services.PutPlay(model.HandModifier{
				MatchId:  handModel.MatchId,
				PlayerId: handModel.PlayerId,
				CardIds:  handModel.Cards,
			})
		}

		handModel.Cards = []int32{}
		return nil
	case " ":
		idx := slices.Index(handModel.HighlightedSlotIdxs, handModel.ActiveSlotIdx)

		if idx > -1 {
			handModel.HighlightedSlotIdxs = slices.Delete(handModel.HighlightedSlotIdxs, 0, 1)
		} else {
			handModel.HighlightedSlotIdxs = append(handModel.HighlightedSlotIdxs, handModel.ActiveSlotIdx)
		}
	case "right":
		handModel.ActiveSlotIdx++

		if handModel.ActiveSlotIdx > len(handModel.Cards)-1 {
			handModel.ActiveSlotIdx = 0
		}

		handModel.HighlighedId = handModel.ActiveSlotIdx //Highlighed id is to be hnalded by view
	case "left":
		handModel.ActiveSlotIdx--

		if handModel.ActiveSlotIdx < 0 {
			handModel.ActiveSlotIdx = len(handModel.Cards) - 1
		}

		handModel.HighlighedId = handModel.ActiveSlotIdx
	}

	return nil
}
func (hc *HandController) Update(msg tea.Msg) tea.Cmd {
	return nil
}
