package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/queries/queries"
	"github.com/charmbracelet/lipgloss"

	"github.com/bardic/gocrib/vo"
)

func GetCardByID(id int, deck *vo.GameDeck) *vo.GameCard {
	for _, card := range deck.Cards {
		gameCard := &vo.GameCard{
			Match: card.Match,
			Card:  card.Card,
		}

		if *card.Match.ID == id {
			return gameCard
		}
	}
	return nil
}

func GetPlayerForAccountID(id *int, match *vo.GameMatch) *vo.GamePlayer {
	for _, player := range match.Players {
		if *player.Accountid == *id {
			return player
		}
	}

	return nil
}

func BuildCommonFooter(match *vo.GameMatch, localplayer *vo.GamePlayer) string {
	currentPlayerTurn := "-"

	if match.Currentplayerturn != nil {
		currentPlayerTurn = strconv.Itoa(*match.Currentplayerturn)
	}

	f := fmt.Sprintf("State: %v\nActive Player: %v\nMatch ID: %d\n", match.Gamestate, currentPlayerTurn, *match.ID)
	playerString := "Players: \n"

	for i, v := range match.Players {
		playerString += styles.PlayerStyles[i].Render(fmt.Sprintf("%d:%d", *v.ID, *v.Score))
		playerString += "\n"
	}

	controls := "tab/shift+tab: navigate screens • space: select • enter: submit • q: exit\n"

	padding := lipgloss.NewStyle().Padding(0, 10, 0, 0)

	cols := lipgloss.JoinHorizontal(lipgloss.Top, padding.Render(f), playerString)

	cols = fmt.Sprintf("\n\n%s\n%s", cols, controls)

	return cols
}

func GetCardSuit(card *queries.Card) string {
	switch card.Suit {
	case queries.CardsuitSpades:
		return "♠"
	case queries.CardsuitHearts:
		return "♥"
	case queries.CardsuitDiamonds:
		return "♦"
	case queries.CardsuitClubs:
		return "♣"
	default:
		return "?"
	}
}

type PegState string

const (
	Filled PegState = "•"
	Empty  PegState = "○"
)

func DrawRow(players []*vo.GamePlayer, pegsToDraw, scoreOffet int) string {
	viewBuilder := strings.Builder{}
	for _, player := range players {
		viewBuilder.WriteString("\n")

		for i := range pegsToDraw {
			if i+scoreOffet == *player.Score {
				viewBuilder.WriteString(string(Filled))
			} else {
				viewBuilder.WriteString(string(Empty))
			}
		}
	}

	viewBuilder.WriteString("\n")
	for range pegsToDraw {
		viewBuilder.WriteString("-")
	}

	return viewBuilder.String()
}

func GetPlayerIDs(players []*vo.GamePlayer) []*int {
	var playIDs []*int
	for _, p := range players {
		playIDs = append(playIDs, p.ID)
	}
	return playIDs
}

func IDFromCards(cards []vo.GameCard) []int {
	var ids []int
	for _, c := range cards {
		ids = append(ids, *c.Match.Cardid)
	}

	return ids
}

func EndPointBuilder(endpoint string, args ...string) string {
	for _, arg := range args {
		endpoint = strings.Replace(endpoint, "%s", arg, 1)
	}
	return endpoint
}
