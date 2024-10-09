package main

import (
	"github.com/bardic/gocrib/cli/styles"
	"github.com/bardic/gocrib/cli/views"
	"github.com/bardic/gocrib/model"
)

type ViewModel struct {
}

func (m *AppModel) View() string {
	var v string
	switch m.ViewStateName {
	case views.Login:
		l := views.LoginView{}
		v = styles.ViewStyle.Render(l.View())
	case views.Lobby:
		l := views.LobbyView{}
		v = styles.ViewStyle.Render(l.View(m.ViewModel))
	case views.Game:
		switch m.GameViewState {
		case model.BoardView:
			v = styles.ViewStyle.Render(views.GameView(
				m.HighlighedId,
				m.HighlightedIds,
				[]model.Card{},
				m.ViewModel,
				m.gameState))
		case model.PlayView:
			v = styles.ViewStyle.Render(views.GameView(
				m.HighlighedId,
				m.HighlightedIds,
				m.play,
				m.ViewModel,
				m.gameState))
		case model.HandView:
			v = styles.ViewStyle.Render(views.GameView(
				m.HighlighedId,
				m.HighlightedIds,
				m.hand,
				m.ViewModel,
				m.gameState))
		case model.KittyView:
			v = styles.ViewStyle.Render(views.GameView(
				m.HighlighedId,
				m.HighlightedIds,
				m.kitty,
				m.ViewModel,
				m.gameState))
		}
	default:
		l := views.LoginView{}
		v = styles.ViewStyle.Render(l.View())
		v = styles.ViewStyle.Render(v)
	}
	return v
}
