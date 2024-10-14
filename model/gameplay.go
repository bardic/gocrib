package model

import (
	"math/rand/v2"
	"time"
)

// _Cards_

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

// _Game_

type GameplayCard struct {
	Card
	CardId    int
	OrigOwner int
	CurrOwner int
	State     CardState
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

type GameMatch struct {
	Id                 int       `json:"id"`
	PlayerIds          []int     `json:"playerIds"`
	PrivateMatch       bool      `json:"privateMatch"`
	EloRangeMin        int       `json:"eloRangeMin"`
	EloRangeMax        int       `json:"eloRangeMax"`
	CreationDate       time.Time `json:"creationDate"`
	DeckId             int       `json:"deckId"`
	CutGameCardId      int       `json:"cutGameCardId"`
	CurrentPlayerTurn  int       `json:"currentPlayerTurn"`
	TurnPassTimestamps []string  `json:"turnPassTimestamps"`
	GameState          GameState `json:"gameState"`
	Art                string    `json:"art "`
	Players            []Player
}

type GameState uint

const (
	NewGameState GameState = 1 << iota
	WaitingState
	MatchReady
	DealState
	CutState
	DiscardState
	PlayState
	OpponentState
	KittyState
	GameWonState
	GameLostState
	MaxGameState
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

type GameViewState uint

const (
	BoardView GameViewState = iota
	PlayView
	HandView
	KittyView
)

// _Player_

type Player struct {
	Id        int
	AccountId int `json:"accountid"`
	Play      []int
	Hand      []int
	Kitty     []int
	Score     int
	Art       string
}

type Account struct {
	Id   int
	Name string
}

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
	GameState GameState
}

type HandModifier struct {
	MatchId  int
	PlayerId int
	CardIds  []int
}

type ScoreResults struct {
	Results []Scores
}

type Scores struct {
	Cards []GameplayCard
	Point int
}

type IViewState interface {
	Enter()
	View() string
}

type ViewStateName uint

const (
	LoginView ViewStateName = iota
	LobbyView
	GameView
	GameOverView
)

var GameTabNames = []string{"Board", "Play", "Hand", "Kitty"}

var LobbyTabNames = []string{"Active Matches", "Open Matches"}

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
