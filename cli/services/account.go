package services

import (
	"net/http"
	"strconv"

	"github.com/bardic/cribbagev2/cli/state"
	tea "github.com/charmbracelet/bubbletea"
)

func Login() tea.Msg {
	return url(EndPointLogin, http.MethodPost, strconv.Itoa(state.PlayerId))
}
