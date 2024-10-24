package views

import (
	"fmt"
	"slices"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/charmbracelet/lipgloss"
)

type KittyView struct {
	HandModel
}

func (v *KittyView) View() string {
	s := ""

	cardViews := make([]string, 0)
	for i := 0; i < len(v.cards); i++ {
		c := utils.GetCardById(v.cards[i], v.deck)
		view := fmt.Sprintf("%v%v", utils.GetCardSuit(c), c.Value)

		if slices.Contains(v.selectedCardIds, c.ID) {
			if i == v.selectedCardId {
				cardViews = append(cardViews, styles.SelectedFocusedStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.SelectedStyle.Render(view))
			}
		} else {
			if i == v.selectedCardId {
				cardViews = append(cardViews, styles.FocusedModelStyle.Render(view))
			} else {
				cardViews = append(cardViews, styles.ModelStyle.Render(view))
			}
		}
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)

	// s += styles.HelpStyle.Render(utils.BuildSubtext(v.player, v.account, utils.IsPlayerTurn(v.player.ID, v.currentTurnPlayerId)))
	return s
}

func (v *KittyView) BuildHeader() string {
	return ""
}

func (v *KittyView) BuildFooter() string {
	return ""
}
