// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package queries

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Cardstate string

const (
	CardstateDeck  Cardstate = "Deck"
	CardstateHand  Cardstate = "Hand"
	CardstatePlay  Cardstate = "Play"
	CardstateKitty Cardstate = "Kitty"
)

func (e *Cardstate) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Cardstate(s)
	case string:
		*e = Cardstate(s)
	default:
		return fmt.Errorf("unsupported scan type for Cardstate: %T", src)
	}
	return nil
}

type NullCardstate struct {
	Cardstate Cardstate
	Valid     bool // Valid is true if Cardstate is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCardstate) Scan(value interface{}) error {
	if value == nil {
		ns.Cardstate, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Cardstate.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCardstate) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Cardstate), nil
}

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
	GamestateNew       Gamestate = "New"
	GamestateWaiting   Gamestate = "Waiting"
	GamestateReady     Gamestate = "Ready"
	GamestateDetermine Gamestate = "Determine"
	GamestateDeal      Gamestate = "Deal"
	GamestateDiscard   Gamestate = "Discard"
	GamestateCut       Gamestate = "Cut"
	GamestatePlay      Gamestate = "Play"
	GamestatePassTurn  Gamestate = "PassTurn"
	GamestateCount     Gamestate = "Count"
	GamestateKitty     Gamestate = "Kitty"
	GamestateWon       Gamestate = "Won"
	GamestateLost      Gamestate = "Lost"
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
	ID   *int
	Name string
}

type Card struct {
	ID    *int
	Value Cardvalue
	Suit  Cardsuit
	Art   string
}

type Deck struct {
	ID             *int
	Cutmatchcardid *int
}

type DeckMatchcard struct {
	Deckid      *int
	Matchcardid *int
}

type GameCard struct {
	Match Matchcard
	Card  Card
}

type Match struct {
	ID                 *int
	Creationdate       pgtype.Timestamptz
	Privatematch       bool
	Elorangemin        *int
	Elorangemax        *int
	Deckid             *int
	Cutgamecardid      *int
	Dealerid           *int
	Currentplayerturn  *int
	Turnpasstimestamps []pgtype.Timestamptz
	Gamestate          Gamestate
	Art                string
}

type MatchPlayer struct {
	Matchid   *int
	Playerid  *int
	Turnorder *int
}

type Matchcard struct {
	ID        *int
	Cardid    *int
	Origowner *int
	Currowner *int
	State     Cardstate
}

type Player struct {
	ID        *int
	Accountid *int
	Score     *int
	Isready   bool
	Art       string
}
