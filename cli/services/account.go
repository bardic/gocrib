package services

import (
	"net/http"

	"github.com/bardic/gocrib/cli/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func Login(accountId string) tea.Msg {
	endPoint := utils.EndPointBuilder(EndPointLogin, accountId)
	return url(endPoint, http.MethodPost, "")
	// return url(EndPointLogin, http.MethodPost, "")
}
