package model

import "time"

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
	Cut
)

type Card struct {
	Id        int `json:"-"`
	Value     CardValue
	Suit      Suit
	OrigOwner int
	CurrOwner int
	State     CardState
	Art       string
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type PlayerHand struct {
	Items []Card
	Size  int
}

type Player struct {
	Id    int `json:"-"`
	Hand  PlayerHand
	Kitty PlayerHand
	Score int
	Art   string
}

type Board struct {
	Players []Player
	Art     string
}

type History struct {
	MatchId               int
	MatchCompletetionDate string
	MatchOutcome          MatchOutcome
}

type MatchOutcome struct {
	Winners []int
	Losers  []int
}

type Chat struct {
	Id       int `json:"-"`
	Members  []int
	Messages []ChatMessage
}

type ChatMessage struct {
	Message   string
	Timestamp string
}

type Lobby struct {
	Id             int `json:"-"`
	Accounts       []int
	CreatationDate string
	PrivateMatch   bool
	EloRangeMin    int
	EloRangeMax    int
}

type Match struct {
	Id                 int `json:"-"`
	LobbyId            int
	BoardId            int
	CurrentPlayerTurn  int
	TurnPassTimestamps []string
}
