package vo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID   int
	Name string
}

type MatchState struct {
	ID    int
	State string
}

type Player struct {
	ID        int
	Accountid int
	Score     int
	Isready   bool
	Art       string
	TurnOrder int
	Hand      []*Card
	Play      []*Card
	Kitty     []*Card
}

type Match struct {
	ID                 int
	Creationdate       pgtype.Timestamptz
	Privatematch       bool
	Elorangemin        int
	Elorangemax        int
	Cutgamecardid      int
	Dealerid           int
	Currentplayerturn  int
	Turnpasstimestamps []pgtype.Timestamptz
	Gamestate          string
	Art                string
	Players            []*Player
}

type HandScore struct {
	ActivePlayerID int
	PlayerSeekID   int
	CardsInPlay    []*Card
	Players        []*Player
}

type Deck struct {
	ID             int
	Cutmatchcardid int
	Matchid        int
	Cards          []*Card
}

type Card struct {
	ID        int
	Cardid    int
	Origowner int
	Currowner int
	State     string
	Name      string
	Rank      int
	Value     int
	Suit      string
	Art       string
}

type HandModifier struct {
	CardIDs []int
}

type ScoreResults struct {
	Results []*Scores
}

type Scores struct {
	Cards []*Card
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
