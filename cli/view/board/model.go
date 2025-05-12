package board

import "github.com/bardic/gocrib/queries/queries"

type Model struct {
	AccountID   *int
	CutIndex    string
	GameMatchID *int
	Gamestate   queries.Gamestate
}
