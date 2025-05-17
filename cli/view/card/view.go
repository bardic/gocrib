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
	LocalPlayer     *vo.Player
	Tabname         string
}

func NewCardView(localPlayer *vo.Player, tabName string) *View {
	return &View{
		LocalPlayer: localPlayer,
		Tabname:     tabName,
	}
}

func (view *View) Update(gameMatch *vo.Match) error {
	return nil
}

func (view *View) Render(gameMatch *vo.Match, gameDeck *vo.Deck, hand []int) string {
	var s string
	var cardViews []string

	s += view.BuildHeader()
	for i := range hand {
		c := utils.GetCardByID(hand[i], gameDeck)

		if c == nil {
			continue
		}

		cardStr := fmt.Sprintf("%v%v", c.Suit, c.Value)
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

		if c.Origowner != view.LocalPlayer.ID {
			styledCard = styles.PlayerStyles[c.Origowner].Render(styledCard)
		} else {
			styledCard = styles.PlayerStyles[c.Currowner].Render(styledCard)
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

func (view *View) BuildFooter(gameMatch *vo.Match) string {
	f := utils.BuildCommonFooter(gameMatch)

	return f
}
