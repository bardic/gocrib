package vo

import "github.com/bardic/gocrib/queries/queries"

type GamePlayer struct {
	*queries.Player
	TurnOrder int
	Hand      []GameCard
	Play      []GameCard
	Kitty     []GameCard
}

type GameMatch struct {
	*queries.Match
	Players []*GamePlayer
}

type ScoreMatch struct {
	ActivePlayerID *int
	PlayerSeekID   *int
	CardsInPlay    *[]GameCard
	Players        *[]GamePlayer
}

type GameDeck struct {
	*queries.Deck
	Cards []*GameCard
}

type GameCard struct {
	Match queries.Matchcard
	Card  queries.Card
}

type GameCardDetails struct {
	Value *int
	Order *int
}

type MatchDetailsResponse struct {
	MatchID   *int
	PlayerID  *int
	GameState queries.Gamestate
}

type HandModifier struct {
	CardIDs []int
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
	AccountID *int
	MatchID   *int
}

type ChangeTabMsg struct {
	TabIndex int
}

type PlayerReady struct {
	MatchID  *int
	PlayerID *int
}

type Hand struct {
	Cards []queries.Matchcard
}
