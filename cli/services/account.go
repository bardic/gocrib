package services

import (
	"net/http"

	"github.com/bardic/gocrib/cli/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func Login(accountID string) tea.Msg {
	endPoint := utils.EndPointBuilder(EndPointLogin, accountID)
	return url(endPoint, http.MethodPost, "")
}
