package services

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func Login(accountId string) tea.Msg {
	return url(EndPointLogin, http.MethodPost, accountId)
}
