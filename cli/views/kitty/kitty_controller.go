package kitty

import (
	"github.com/bardic/gocrib/cli/views"
	tea "github.com/charmbracelet/bubbletea"
)

type KittyController struct {
	views.Controller
}

func (gc *KittyController) GetState() views.ControllerState {
	return views.LobbyControllerState
}

func (gc *KittyController) Init() {
	gc.Model = KittyModel{}
}
func (gc *KittyController) Render() string {
	return ""
}

func (hc *KittyController) InitView() {

}
func (hc *KittyController) ParseInput(msg tea.KeyMsg) tea.Msg {
	kittyModel := hc.Model.(*KittyModel)
	switch msg.String() {
	case "right":
		kittyModel.ActiveSlotIdx++

		if kittyModel.ActiveSlotIdx > len(kittyModel.Cards)-1 {
			kittyModel.ActiveSlotIdx = 0
		}

		kittyModel.HighlighedId = kittyModel.ActiveSlotIdx //Highlighed id is to be hnalded by view
	case "left":
		kittyModel.ActiveSlotIdx--

		if kittyModel.ActiveSlotIdx < 0 {
			kittyModel.ActiveSlotIdx = len(kittyModel.Cards) - 1
		}

		kittyModel.HighlighedId = kittyModel.ActiveSlotIdx
	}

	return nil
}
func (hc *KittyController) Update(msg tea.Msg) tea.Cmd {
	return nil
}
