package vo

import "github.com/bardic/gocrib/queries/queries"

type GameMatch struct {
	queries.Match
	Players []*queries.Player
}

type GameDeck struct {
	*queries.Deck
	Cards []*GameCard
}

type GameCard struct {
	queries.Matchcard
	queries.Card
}

type GameCardDetails struct {
	Value *int
	Order *int
}

type GameActionType int

const (
	Cut GameActionType = iota
	Discard
	Peg
	Tally
)

type GameAction struct {
	MatchId  *int
	Type     GameActionType
	CardsIds []*int
}

//_Comms_

type MatchRequirements struct {
	AccountId   *int
	IsPrivate   bool
	EloRangeMin *int
	EloRangeMax *int
}

type CutDeckReq struct {
	MatchId  *int
	CutIndex string
}

type MatchDetailsResponse struct {
	MatchId   *int
	PlayerId  *int
	GameState queries.Gamestate
}

type HandModifier struct {
	SenderId *int
	CardIds  []*int
}

type ScoreResults struct {
	Results []Scores
}

type Scores struct {
	Cards []GameCard
	Point *int
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
	NewState  ViewStateName
	AccountId *int
	MatchId   *int
}

type GameStateChangeMsg struct {
	NewState queries.Gamestate
	PlayerId *int
	MatchId  *int
}

type ChangeTabMsg struct {
	TabIndex *int
}

type UIFooterVO struct {
	ActivePlayerId *int
	MatchId        *int
	GameState      queries.Gamestate
	LocalPlayerID  *int
}

type PlayerReady struct {
	MatchId  *int // MatchId
	PlayerId *int // PlayerId
}

type Kitty struct {
	Cards []*int
}
