package model

import (
	"math/rand/v2"
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

func (d *GameDeck) Shuffle() *GameDeck {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})

	return d
}

type Player struct {
	Id        int32
	AccountId int
	Play      []int
	Hand      []int
	Kitty     []int
	Score     int
	Art       string
}

func (p *Player) Eq(c Player) bool {
	if p.Id != c.Id {
		return false
	}

	if p.AccountId != c.AccountId {
		return false
	}

	if p.Score != c.Score {
		return false
	}

	if p.Art != c.Art {
		return false
	}

	if !eqIntArr(p.Hand, c.Hand) {
		return false
	}

	if !eqIntArr(p.Play, c.Play) {
		return false
	}

	if !eqIntArr(p.Kitty, c.Kitty) {
		return false
	}

	return true
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
	RequesterId int32
	IsPrivate   bool
	EloRangeMin int
	EloRangeMax int
}

type CardSlots uint

const (
	CardOne CardSlots = iota
	CardTwo
	CardThree
	CardFour
	CardFive
	CardSix
)

type MatchDiffType uint

const (
	GenericDiff = 1 << iota
	NewDeckDiff
	CutDiff
	TurnDiff
	GameStateDiff
	CardsInPlayDiff
	TurnPassTimestampsDiff
	MaxDiff
)

type GameMatch struct {
	Match
	Players []Player
}

type Match struct {
	Id                 int       `json:"id"`
	PlayerIds          []int     `json:"playerIds"`
	PrivateMatch       bool      `json:"privateMatch"`
	EloRangeMin        int       `json:"eloRangeMin"`
	EloRangeMax        int       `json:"eloRangeMax"`
	CreationDate       time.Time `json:"creationDate"`
	DeckId             int       `json:"deckId"`
	CardsInPlay        []int     `json:"cardsInPlay"`
	CutGameCardId      int       `json:"cutGameCardId"`
	CurrentPlayerTurn  int       `json:"currentPlayerTurn"`
	TurnPassTimestamps []string  `json:"turnPassTimestamps"`
	GameState          GameState `json:"gameState"`
	Art                string    `json:"art "`
}

func (m *Match) Eq(c Match) int {

	diff := 0

	if m.Id != c.Id {
		diff |= GenericDiff
	}

	if m.PrivateMatch != c.PrivateMatch {
		diff |= GenericDiff
	}

	if m.EloRangeMin != c.EloRangeMin {
		diff |= GenericDiff
	}

	if m.EloRangeMax != c.EloRangeMax {
		diff |= GenericDiff
	}

	if m.CreationDate != c.CreationDate {
		diff |= GenericDiff
	}

	if m.DeckId != c.DeckId {
		diff |= NewDeckDiff
	}

	if m.CutGameCardId != c.CutGameCardId {
		diff |= CutDiff
	}

	if m.CurrentPlayerTurn != c.CurrentPlayerTurn {
		diff |= TurnDiff
	}

	if m.GameState != c.GameState {
		diff |= GameStateDiff
	}

	if m.Art != c.Art {
		diff |= GenericDiff
	}

	if !eqIntArr(m.PlayerIds, c.PlayerIds) {
		diff |= GenericDiff
	}

	if !eqIntArr(m.CardsInPlay, c.CardsInPlay) {
		diff |= CardsInPlayDiff
	}

	if !eqStringArr(m.TurnPassTimestamps, c.TurnPassTimestamps) {
		diff |= TurnPassTimestampsDiff
	}

	return diff
}

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

func eqStringArr(a, b []string) bool {
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

func eqPlayerArr(a, b []Player) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].Eq(b[i]) {
			return false
		}
	}
	return true
}

type GameState uint

const (
	WaitingState GameState = 1 << iota
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
	MatchId  int
	Type     GameActionType
	CardsIds []int
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
	ActiveView ViewState = iota
	LobbyView
	BoardView
	PlayView
	HandView
	KittyView
	ScoresView
	GameOverView
)

var TabNames = []string{"Board", "Play", "Hand", "Kitty"}

type HandModifier struct {
	MatchId  int
	PlayerId int
	CardIds  []int
}
