package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"
	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayerByForMatchAndAccount(matchId, playerId *int) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointPlayerByForMatchAndAccount, strconv.Itoa(*matchId), strconv.Itoa(*playerId))

	return url(endpoint, http.MethodGet, "")
}

func PutKitty(matchId, fromPlayerId, toPlayerId *int, cards vo.HandModifier) tea.Msg {
	b, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	endpoint := utils.EndPointBuilder(EndPointKitty, strconv.Itoa(*matchId), strconv.Itoa(*fromPlayerId), strconv.Itoa(*toPlayerId))

	return url(endpoint, http.MethodPut, string(b))
}

func PutPlay(matchId, fromPlayerId, toPlayerId *int, cards vo.HandModifier) tea.Msg {
	b, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	endpoint := utils.EndPointBuilder(EndPointPlay, strconv.Itoa(*matchId), strconv.Itoa(*fromPlayerId), strconv.Itoa(*toPlayerId))

	return url(endpoint, http.MethodPut, string(b))
}
