package card

import (
	"fmt"
	"slices"

	"cli/styles"
	"cli/utils"
	cliVO "cli/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	*cliVO.HandModel
	SelectedCardId int
}

func (view *View) Init() {
	view.SelectedCardId = 0
}

func (view *View) Render() string {
	s := ""
	cardViews := make([]string, 0)

	for i := 0; i < len(view.CardsToDisplay); i++ {
		c := utils.GetCardById(view.CardsToDisplay[i], view.Deck)
		cardStr := fmt.Sprintf("%v%v", utils.GetCardSuit(c), c.Value)

		if slices.Contains(view.SelectedCardIds, c.ID) {
			if i == view.SelectedCardId {
				cardViews = append(cardViews, styles.SelectedFocusedStyle.Render(cardStr))
			} else {
				cardViews = append(cardViews, styles.SelectedStyle.Render(cardStr))
			}
		} else {
			if i == view.SelectedCardId {
				cardViews = append(cardViews, styles.FocusedModelStyle.Render(cardStr))
			} else {
				cardViews = append(cardViews, styles.ModelStyle.Render(cardStr))
			}
		}
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)
	return s
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	return ""
}
