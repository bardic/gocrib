package model

import (
	"time"
)

type CardValue int

const (
	Ace CardValue = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Joker
)

type Suit int

const (
	Spades Suit = iota
	Clubs
	Hearts
	Diamonds
)

type CardState int

const (
	Deck CardState = iota
	Hand
	Play
	Kitty
)

type Card struct {
	Id    int       `json:"Id"`
	Value CardValue `json:"Value"`
	Suit  Suit      `json:"Suit"`
	Art   string    `json:"Art"`
}

type Cards struct {
	Cards []Card `json:"Cards"`
}

type GameplayCard struct {
	Id        int
	CardId    int
	OrigOwner int
	CurrOwner int
	State     CardState
	Value     int    `json:"-"`
	Suit      string `json:"-"`
	Art       string `json:"-"`
}

type GameDeck struct {
	Id    int
	Cards []GameplayCard `json:"cards"`
}

type Player struct {
	Id        int
	AccountId int
	Hand      []int
	Kitty     []int
	Score     int
	Art       string
}

type History struct {
	MatchId               int
	MatchCompletetionDate string
	Winners               []int
	Losers                []int
}

type Chat struct {
	Id       int `json:"-"`
	Members  []int
	Messages []ChatMessage
}

type ChatMessage struct {
	Sender    int
	Message   string
	Timestamp string
}

type MatchRequirements struct {
	IsPrivate   bool
	EloRangeMin int
	EloRangeMax int
}

type Match struct {
	Id                 int
	PlayerIds          []int
	PrivateMatch       bool
	EloRangeMin        int
	EloRangeMax        int
	CreationDate       time.Time
	DeckId             int
	CardsInPlay        []int
	CutGameCardId      int
	CurrentPlayerTurn  int
	TurnPassTimestamps []string
	GameState          GameState
	Art                string
	Players            []Player
}

type GameState uint

const (
	WaitingState GameState = iota
	DealState
	CutState
	DiscardState
	PlayState
	OpponentState
	KittyState
	GameWonState
	GameLostState
)

type GameActionType int

const (
	Cut GameActionType = iota
	Discard
	Peg
	Tally
)

type GameAction struct {
	MatchId        int
	Type           GameActionType
	GameplayCardId int
}

type ScoreResults struct {
	Results []Scores
}

type Scores struct {
	Cards []GameplayCard
	Point int
}

type ViewState uint

const (
	LobbyView ViewState = iota
	KittyView
	PlayView
	HandView
	ScoresView
	GameOverView
)
