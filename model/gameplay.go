package model

import (
	"github.com/bardic/gocrib/queries"
)

// _Cards_

// type CardValue int

// const (
// 	Ace CardValue = iota
// 	Two
// 	Three
// 	Four
// 	Five
// 	Six
// 	Seven
// 	Eight
// 	Nine
// 	Ten
// 	Jack
// 	Queen
// 	King
// 	Joker
// )

// type Suit int

// const (
// 	Spades Suit = iota
// 	Clubs
// 	Hearts
// 	Diamonds
// )

// type CardState int

// const (
// 	Deck CardState = iota
// 	Hand
// 	Play
// 	Kitty
// )

// type Card struct {
// 	Id    int       `json:"Id"`
// 	Value CardValue `json:"Value"`
// 	Suit  Suit      `json:"Suit"`
// 	Art   string    `json:"Art"`
// }

type Cards struct {
	Cards []queries.Card `json:"Cards"`
}

// _Game_

// type GameplayCard struct {
// 	Card
// 	CardId    int
// 	OrigOwner int
// 	CurrOwner int
// 	State     CardState
// }

// type GameDeck struct {
// 	Id    int
// 	Cards []GameplayCard `json:"cards"`
// }

// type GameMatch struct {
// 	Id                 int       `json:"id"`
// 	PlayerIds          []int     `json:"playerIds"`
// 	PrivateMatch       bool      `json:"privateMatch"`
// 	EloRangeMin        int       `json:"eloRangeMin"`
// 	EloRangeMax        int       `json:"eloRangeMax"`
// 	CreationDate       time.Time `json:"creationDate"`
// 	DeckId             int       `json:"deckId"`
// 	CutGameCardId      int       `json:"cutGameCardId"`
// 	CurrentPlayerTurn  int       `json:"currentPlayerTurn"`
// 	TurnPassTimestamps []string  `json:"turnPassTimestamps"`
// 	GameState          GameState `json:"gameState"`
// 	Art                string    `json:"art "`
// 	Players            []Player
// }

type GameMatch struct {
	queries.Match
	Players []queries.Player
}

type GameDeck struct {
	queries.Deck
	Cards []queries.Gameplaycard
}

// type GameState uint

// const (
// 	NewGameState GameState = 1 << iota
// 	JoinGameState
// 	WaitingState
// 	MatchReady
// 	DealState
// 	DiscardState
// 	CutState
// 	PlayState
// 	OpponentState
// 	KittyState
// 	GameWonState
// 	GameLostState
// 	MaxGameState
// )

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

type GameViewState uint

const (
	BoardView GameViewState = iota
	PlayView
	HandView
	KittyView
)

// _Player_

// type Player struct {
// 	Id        int
// 	AccountId int `json:"accountid"`
// 	Play      []int
// 	Hand      []int
// 	Kitty     []int
// 	Score     int
// 	IsReady   bool
// 	Art       string
// }

//_Comms_

type MatchRequirements struct {
	PlayerId    int
	AccountId   int
	IsPrivate   bool
	EloRangeMin int
	EloRangeMax int
}

type CutDeckReq struct {
	PlayerId int
	MatchId  int
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
	Cards []queries.Gameplaycard
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

type LobbyViewState uint

const (
	OpenMatches LobbyViewState = iota
	AvailableMatches
)

func eqIntArr(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type StateChangeMsg struct {
	NewState ViewStateName
	MatchId  int
}

type GameStateChangeMsg struct {
	NewState queries.Gamestate
	PlayerId int
	MatchId  int
}
