package core

import (
	"github.com/bardic/gocrib/queries/queries"
)

type State struct {
	state   *queries.Gamestate
	matchID *int
}

func (s *State) GetState() *queries.Gamestate {
	return s.state
}

func (s *State) SetState(state *queries.Gamestate) error {
	s.state = state
	return nil
}
