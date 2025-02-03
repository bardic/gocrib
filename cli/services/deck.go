package services

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func GetDeckByMatchId(id int) tea.Msg {
	endPoint := utils.EndPointBuilder(EndPointDeckById, strconv.Itoa(id))
	return url(endPoint, http.MethodGet, "")
}

func GetDeckByPlayIdAndMatchId(playerId, matchId int) tea.Msg {
	endPoint := utils.EndPointBuilder(EndPointDeckByPlayerAndMatchId, strconv.Itoa(matchId), strconv.Itoa(playerId))
	return url(endPoint, http.MethodGet, "")
}
