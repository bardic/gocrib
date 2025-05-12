package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"
	"github.com/bardic/gocrib/vo"
	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayerByForMatchAndAccount(matchID, playerID *int) tea.Msg {
	endpoint := utils.EndPointBuilder(
		EndPointPlayerByForMatchAndAccount,
		strconv.Itoa(*matchID),
		strconv.Itoa(*playerID),
	)

	return url(endpoint, http.MethodGet, "")
}

func PutKitty(matchID, fromPlayerID, toPlayerID *int, cards vo.HandModifier) tea.Msg {
	b, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	endpoint := utils.EndPointBuilder(
		EndPointKitty,
		strconv.Itoa(*matchID),
		strconv.Itoa(*fromPlayerID),
		strconv.Itoa(*toPlayerID),
	)

	return url(endpoint, http.MethodPut, string(b))
}

func PutPlay(matchID, fromPlayerID, toPlayerID *int, cards vo.HandModifier) tea.Msg {
	b, err := json.Marshal(cards)
	if err != nil {
		return err
	}

	endpoint := utils.EndPointBuilder(
		EndPointPlay,
		strconv.Itoa(*matchID),
		strconv.Itoa(*fromPlayerID),
		strconv.Itoa(*toPlayerID),
	)

	return url(endpoint, http.MethodPut, string(b))
}
