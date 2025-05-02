package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"

	tea "github.com/charmbracelet/bubbletea"
)

func GetMatchById(matchId *int) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointMatch, strconv.Itoa(*matchId))
	return url(endpoint, http.MethodGet, "")
}

func GetOpenMatches() tea.Msg {
	return url(EndPointOpenMatch, http.MethodGet, "")
}

func JoinMatch(accountId, activeMatchId int) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointJoinMatch, strconv.Itoa(activeMatchId), strconv.Itoa(accountId))
	return url(endpoint, http.MethodPut, "")
}

func PostPlayerMatch(accountId *int) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointMatch, strconv.Itoa(*accountId))
	return sendReq(endpoint, http.MethodPost, nil)
}

func CutDeck(matchId int, cutIndex string) tea.Msg {
	endpoint := utils.EndPointBuilder(EndPointMatchCutDeck, strconv.Itoa(matchId), cutIndex)

	return sendReq(endpoint, http.MethodPut, "")
}

func sendReq(endPoint string, method string, body interface{}) tea.Msg {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	return url(endPoint, method, string(b))
}
