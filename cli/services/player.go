package services

import (
	"encoding/json"
	"net/http"

	"github.com/bardic/gocrib/cli/state"
	"github.com/bardic/gocrib/model"
	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayer() tea.Msg {
	return url(EndPointPlayer, http.MethodGet, "")
}

func PutPlayer(id int) tea.Msg {
	return url(EndPointPlayer, http.MethodPut, "")
}

func PutKitty() tea.Msg {
	b, err := json.Marshal(state.CurrentHandModifier)
	if err != nil {
		return err
	}

	return url(EndPointKitty, http.MethodPut, string(b))
}

func PutPlay() tea.Msg {
	b, err := json.Marshal(state.CurrentHandModifier)
	if err != nil {
		return err
	}

	return url(EndPointPlay, http.MethodPut, string(b))
}

func PostPlayer(player model.Player) tea.Msg {

	b, err := json.Marshal(player)
	if err != nil {
		return err
	}

	return url(EndPointPlayer, http.MethodPost, string(b))
}

func DeletePlayer(ids []int) tea.Msg {
	return url(EndPointPlayer, http.MethodDelete, "")
}
