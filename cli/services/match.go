package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func GetMatchByID(matchID *int) tea.Msg {
	matchIDStr := ""
	if matchID != nil {
		matchIDStr = strconv.Itoa(*matchID)
	}

	endpoint := utils.EndPointBuilder(EndPointMatch, matchIDStr)
	return url(endpoint, http.MethodGet, "")
}

func GetOpenMatches() tea.Msg {
	return url(EndPointOpenMatch, http.MethodGet, "")
}

func JoinMatch(accountID, activeMatchID int) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointJoinMatch, strconv.Itoa(activeMatchID), strconv.Itoa(accountID))
	return url(endpoint, http.MethodPut, "")
}

func PostPlayerMatch(accountID *int) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointMatch, strconv.Itoa(*accountID))
	return sendReq(endpoint, http.MethodPost, nil)
}

func CutDeck(matchID int, cutIndex string) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointMatchCutDeck, strconv.Itoa(matchID), cutIndex)

	return sendReq(endpoint, http.MethodPut, "")
}

func sendReq(endPoint string, method string, body any) tea.Msg {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return url(endPoint, method, string(b))
}
