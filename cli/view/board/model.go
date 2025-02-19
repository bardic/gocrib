package board

import "github.com/bardic/gocrib/queries/queries"

type Model struct {
	AccountId   *int
	CutIndex    string
	GameMatchId *int
	Gamestate   queries.Gamestate
}
