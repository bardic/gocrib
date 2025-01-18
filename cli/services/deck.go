package services

import (
	"net/http"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

func GetDeckById(id int) tea.Msg {
	return url(EndPointDeckById+strconv.Itoa(int(id)), http.MethodGet, "")
}
