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
	AccountId   int32
	IsPrivate   bool
	EloRangeMin int32
	EloRangeMax int32
}

type CutDeckReq struct {
	MatchId  int32
	CutIndex string
}

type MatchDetailsResponse struct {
	MatchId   int32
	PlayerId  int32
	GameState queries.Gamestate
}

type HandModifier struct {
	CardIds []int32
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
	NewState  ViewStateName
	AccountId int32
	MatchId   int32
}

type GameStateChangeMsg struct {
	NewState queries.Gamestate
	PlayerId int32
	MatchId  int32
}

type ChangeTabMsg struct {
	TabIndex int
}

type UIFooterVO struct {
	ActivePlayerId int32
	MatchId        int32
	GameState      queries.Gamestate
	LocalPlayerID  int32
}

type PlayerReady struct {
	MatchId  int32 // MatchId
	PlayerId int32 // PlayerId
}

type Kitty struct {
	Cards []int32
}
