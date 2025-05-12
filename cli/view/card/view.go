package card

import (
	"fmt"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"

	"github.com/charmbracelet/lipgloss"
)

type View struct {
	ActiveCardID    int
	SelectedCardIDs []int
	LocalPlayer     *vo.GamePlayer
	Tabname         string
}

func NewCardView(localPlayer *vo.GamePlayer, tabName string) *View {
	return &View{
		LocalPlayer: localPlayer,
		Tabname:     tabName,
	}
}

func (view *View) Update(gameMatch *vo.GameMatch) error {
	return nil
}

func (view *View) Render(gameMatch *vo.GameMatch, gameDeck *vo.GameDeck, hand []int) string {
	var s string
	var cardViews []string

	s += view.BuildHeader()
	for i := range hand {
		c := utils.GetCardByID(hand[i], gameDeck)

		if c == nil {
			continue
		}

		cardStr := fmt.Sprintf("%v%v", utils.GetCardSuit(&c.Card), c.Card.Value)
		var styledCard string

		top := lipgloss.Place(0, 0, lipgloss.Right, lipgloss.Bottom, cardStr)

		bottom := lipgloss.Place(8, 5, lipgloss.Right, lipgloss.Bottom, cardStr)

		styledCard = styles.ModelStyle.Render(top, bottom)

		if i == view.ActiveCardID {
			styledCard = styles.SelectedFocusedStyle.Render(top, bottom)
		}

		for _, v := range view.SelectedCardIDs {
			if v == hand[i] {
				styledCard = styles.FocusedModelStyle.Render(top, bottom)
			}
		}

		// Original owner seems to be account id and not match player id

		if c.Match.Origowner != nil {
			styledCard = styles.PlayerStyles[*c.Match.Origowner].Render(styledCard)
		} else {
			styledCard = styles.PlayerStyles[*c.Match.Currowner].Render(styledCard)
		}

		cardViews = append(cardViews, styledCard)
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, cardViews...)

	s += "\n" + view.BuildFooter(gameMatch)

	return s
}

func (view *View) BuildHeader() string {
	return ""
}

func (view *View) BuildFooter(gameMatch *vo.GameMatch) string {
	f := utils.BuildCommonFooter(gameMatch, view.LocalPlayer)

	return f
}
