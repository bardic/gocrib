package card

import (
	"fmt"
	"queries"
	"slices"

	"cli/styles"
	"cli/utils"
	cliVO "cli/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	*cliVO.HandVO
	ActiveCardId    int32
	SelectedCardIds []int32
	ActivePlayerId  int32
	MatchId         int32
	GameState       queries.Gamestate
}

func (view *View) Init() {
	view.ActiveCardId = 0
}

func (view *View) Render() string {
	var s string
	var cardViews []string

	s += view.BuildHeader()
	for i := 0; i < len(view.CardIds); i++ {
		c := utils.GetCardById(view.CardIds[i], view.Deck)
		cardStr := fmt.Sprintf("%v%v", utils.GetCardSuit(c), c.Value)
		styledCard := styles.ModelStyle.Render(cardStr)
		if slices.Index(view.SelectedCardIds, c.ID) > -1 {
			if int32(i) == view.ActiveCardId {
				styledCard = styles.SelectedFocusedStyle.Render(cardStr)
			} else {
				styledCard = styles.FocusedModelStyle.Render(cardStr)
			}
		} else {
			if int32(i) == view.ActiveCardId {
				styledCard = styles.FocusedModelStyle.Render(cardStr)
			}
		}
		cardViews = append(cardViews, styledCard)
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)

	s += "\n" + view.BuildFooter()

	return s
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter() string {
	f := utils.BuildCommonFooter(
		view.ActivePlayerId,
		view.LocalPlayerID,
		view.MatchId,
		view.GameState,
	)

	return f
}
