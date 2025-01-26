package services

import (
	"net/http"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func GetDeckByMatchId(id int) tea.Msg {
	endPoint := endPointBuilder(EndPointDeckById, strconv.Itoa(id))
	return url(endPoint, http.MethodGet, "")
}

func endPointBuilder(endpoint string, args ...string) string {
	for _, arg := range args {
		endpoint = strings.Replace(endpoint, "%s", arg, 1)
	}
	return endpoint
}
