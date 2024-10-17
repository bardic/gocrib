// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Cardsuit string

const (
	CardsuitSpades   Cardsuit = "Spades"
	CardsuitClubs    Cardsuit = "Clubs"
	CardsuitHearts   Cardsuit = "Hearts"
	CardsuitDiamonds Cardsuit = "Diamonds"
)

func (e *Cardsuit) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Cardsuit(s)
	case string:
		*e = Cardsuit(s)
	default:
		return fmt.Errorf("unsupported scan type for Cardsuit: %T", src)
	}
	return nil
}

type NullCardsuit struct {
	Cardsuit Cardsuit
	Valid    bool // Valid is true if Cardsuit is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCardsuit) Scan(value interface{}) error {
	if value == nil {
		ns.Cardsuit, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Cardsuit.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCardsuit) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Cardsuit), nil
}

type Cardvalue string

const (
	CardvalueAce   Cardvalue = "Ace"
	CardvalueTwo   Cardvalue = "Two"
	CardvalueThree Cardvalue = "Three"
	CardvalueFour  Cardvalue = "Four"
	CardvalueFive  Cardvalue = "Five"
	CardvalueSix   Cardvalue = "Six"
	CardvalueSeven Cardvalue = "Seven"
	CardvalueEight Cardvalue = "Eight"
	CardvalueNine  Cardvalue = "Nine"
	CardvalueTen   Cardvalue = "Ten"
	CardvalueJack  Cardvalue = "Jack"
	CardvalueQueen Cardvalue = "Queen"
	CardvalueKing  Cardvalue = "King"
	CardvalueJoker Cardvalue = "Joker"
)

func (e *Cardvalue) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Cardvalue(s)
	case string:
		*e = Cardvalue(s)
	default:
		return fmt.Errorf("unsupported scan type for Cardvalue: %T", src)
	}
	return nil
}

type NullCardvalue struct {
	Cardvalue Cardvalue
	Valid     bool // Valid is true if Cardvalue is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCardvalue) Scan(value interface{}) error {
	if value == nil {
		ns.Cardvalue, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Cardvalue.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCardvalue) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Cardvalue), nil
}

type Gamestate string

const (
	GamestateNewGameState  Gamestate = "NewGameState"
	GamestateJoinGameState Gamestate = "JoinGameState"
	GamestateWaitingState  Gamestate = "WaitingState"
	GamestateMatchReady    Gamestate = "MatchReady"
	GamestateDealState     Gamestate = "DealState"
	GamestateDiscardState  Gamestate = "DiscardState"
	GamestateCutState      Gamestate = "CutState"
	GamestatePlayState     Gamestate = "PlayState"
	GamestateOpponentState Gamestate = "OpponentState"
	GamestateKittyState    Gamestate = "KittyState"
	GamestateGameWonState  Gamestate = "GameWonState"
	GamestateGameLostState Gamestate = "GameLostState"
	GamestateMaxGameState  Gamestate = "MaxGameState"
)

func (e *Gamestate) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Gamestate(s)
	case string:
		*e = Gamestate(s)
	default:
		return fmt.Errorf("unsupported scan type for Gamestate: %T", src)
	}
	return nil
}

type NullGamestate struct {
	Gamestate Gamestate
	Valid     bool // Valid is true if Gamestate is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullGamestate) Scan(value interface{}) error {
	if value == nil {
		ns.Gamestate, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Gamestate.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullGamestate) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Gamestate), nil
}

type Account struct {
	ID   int32
	Name string
}

type Card struct {
	ID    int32
	Value Cardvalue
	Suit  Cardsuit
	Art   string
}

type Chat struct {
	ID       int32
	Members  []int32
	Messages []int32
}

type Chatmessage struct {
	ID        int32
	Sender    int32
	Message   string
	Timestamp pgtype.Timestamptz
}

type Deck struct {
	ID    int32
	Cards []byte
}

type Gameplaycard struct {
	ID        int32
	Cardid    int32
	Origowner pgtype.Int4
	Currowner pgtype.Int4
	State     Gamestate
}

type Match struct {
	ID                 int32
	Playerids          []int32
	Creationdate       pgtype.Timestamptz
	Privatematch       bool
	Elorangemin        int32
	Elorangemax        int32
	Deckid             int32
	Cutgamecardid      int32
	Currentplayerturn  int32
	Turnpasstimestamps []pgtype.Timestamptz
	Gamestate          int32
	Art                string
}

type Matchhistory struct {
	ID                    int32
	Matchid               int32
	Matchcompletetiondate pgtype.Timestamptz
	Winners               []int32
	Losers                []int32
}

type Player struct {
	ID        int32
	Accountid int32
	Play      []int32
	Hand      []int32
	Kitty     []int32
	Score     int32
	Isready   bool
	Art       string
}