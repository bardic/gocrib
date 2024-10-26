package vo

import (
	"queries"
)

type GameMatch struct {
	queries.Match
	Players []*queries.Player
}

type GameDeck struct {
	*queries.Deck
	Cards []queries.GetGameCardsForMatchRow
}

type GameCard struct {
	queries.Matchcard
	queries.Card
}

type GameCardDetails struct {
	Value int
	Order int
}

type GameActionType int

const (
	Cut GameActionType = iota
	Discard
	Peg
	Tally
)

type GameAction struct {
	MatchId  int
	Type     GameActionType
	CardsIds []int
}

//_Comms_

type MatchRequirements struct {
	PlayerId    int32
	AccountId   int32
	IsPrivate   bool
	EloRangeMin int32
	EloRangeMax int32
}

type CutDeckReq struct {
	PlayerId int32
	MatchId  int32
	CutIndex string
}

type JoinMatchReq struct {
	MatchId  int
	PlayerId int
}

type MatchDetailsResponse struct {
	MatchId   int
	PlayerId  int
	GameState queries.Gamestate
}

type HandModifier struct {
	MatchId  int32
	PlayerId int32
	CardIds  []int32
}

type ScoreResults struct {
	Results []Scores
}

type Scores struct {
	Cards []GameCard
	Point int
}

type ViewStateName uint

const (
	LoginView ViewStateName = iota
	LobbyView
	CreateGameView
	JoinGameView
	PlayersReadyView
	InGameView
	GameOverView
)

type ViewState uint

const (
	OpenMatches ViewState = iota
	AvailableMatches
	BoardView
	PlayView
	HandView
	KittyView
)

type StateChangeMsg struct {
	NewState ViewStateName
	MatchId  int
}

type GameStateChangeMsg struct {
	NewState queries.Gamestate
	PlayerId int32
	MatchId  int32
}

type ChangeTabMsg struct {
	TabIndex int
}
