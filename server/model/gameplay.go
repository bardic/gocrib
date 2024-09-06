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
	Id    int
	Value CardValue
	Suit  Suit
	Art   string
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

type Player struct {
	Id    int
	Hand  []int
	Kitty []int
	Score int
	Art   string
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

type Lobby struct {
	Id             int
	AccountIds     []int
	PrivateMatch   bool
	EloRangeMin    int
	EloRangeMax    int
	CreatationDate time.Time
}

type Match struct {
	Id                 int
	LobbyId            int
	DeckId             int
	PlayerIds          []int
	CardsInPlay        []int
	CutGameCardId      int
	CurrentPlayerTurn  int
	TurnPassTimestamps []string
	Art                string
	Players            []Player
}

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

type GameDeck struct {
	Id      int
	MatchId int
	Cards   []int //GameplayCardId
}
