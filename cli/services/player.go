package services

import (
	"encoding/json"
	"net/http"

	"github.com/bardic/gocrib/model"
	tea "github.com/charmbracelet/bubbletea"
)

func GetPlayer() tea.Msg {
	return url(EndPointPlayer, http.MethodGet, "")
}

func PutPlayer(id int) tea.Msg {
	return url(EndPointPlayer, http.MethodPut, "")
}

func PutKitty(req model.HandModifier) tea.Msg {
	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return url(EndPointKitty, http.MethodPut, string(b))
}

func PutPlay(req model.HandModifier) tea.Msg {
	b, err := json.Marshal(req)
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
