package player

import (
	"testing"
)

var (
	playerId  = 1
	accountId = 1
	score     = 0
	// player    = vo.GamePlayer{
	// 	Player: queries.Player{
	// 		ID:        &playerId,
	// 		Accountid: &accountId,
	// 		Score:     &score,
	// 		Isready:   true,
	// 	},
	// 	TurnOrder: 1,
	// 	Hand:      []vo.GameCard{},
	// 	Play:      []vo.GameCard{},
	// 	Kitty:     []vo.GameCard{},
	// }
)

// func TestGetPlayer(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	m := NewMockFoo(ctrl)

// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodGet, "/player/1", nil)
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Assertions
// 	if assert.NoError(t, GetPlayer(c)) {
// 		assert.Equal(t, http.StatusCreated, rec.Code)
// 		assert.Equal(t, player, rec.Body.String())
// 	}
// }

func TestBuildGetPlayer(t *testing.T) {
	// ctrl := gomock.NewController(t)
	// m := NewMockFoo(ctrl)

	// e := echo.New()
	// req := httptest.NewRequest(http.MethodGet, "/player/1", nil)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// rec := httptest.NewRecorder()
	// c := e.NewContext(req, rec)

	// // Assertions
	// if assert.NoError(t, GetPlayer(c)) {
	// 	assert.Equal(t, http.StatusCreated, rec.Code)
	// 	assert.Equal(t, player, rec.Body.String())
	// }
}

// func Test_createGameCard(t *testing.T) {
// 	type args struct {
// 		deck  vo.GameDeck
// 		cards []queries.Matchcard
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want []vo.GameCard
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := createGameCard(tt.args.deck, tt.args.cards); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("createGameCard() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
