package services

import (
	"net/http"
	"strconv"

	"github.com/bardic/gocrib/cli/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func GetDeckByPlayIDAndMatchID(playerID, matchID int) tea.Msg {
	endPoint := utils.EndPointBuilder(EndPointDeckByPlayerAndMatchID, strconv.Itoa(matchID), strconv.Itoa(playerID))
	return url(endPoint, http.MethodGet, "")
}
