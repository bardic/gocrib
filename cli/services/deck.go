package services

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func GetDeckByPlayIdAndMatchId(playerId, matchId int) tea.Msg {
	endPoint := utils.EndPointBuilder(EndPointDeckByPlayerAndMatchId, strconv.Itoa(matchId), strconv.Itoa(playerId))
	return url(endPoint, http.MethodGet, "")
}
