package services

import (
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

func GetHistory() tea.Msg {
	return url(EndPointHistory, http.MethodGet, "")

}

func PutHistory(id int) tea.Msg {
	return url(EndPointHistory, http.MethodPut, "")
}

func PostHistory(ids []int) tea.Msg {
	return url(EndPointHistory, http.MethodPost, "")
}

func DeleteHistory(ids []int) tea.Msg {
	return url(EndPointHistory, http.MethodDelete, "")
}
