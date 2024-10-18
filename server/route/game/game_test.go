package game

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/bardic/gocrib/model"
// )

// // TestHelloName calls greetings.Hello with a name, checking
// // for a valid return value.
// func TestThirtyOne(t *testing.T) {
// 	hand := []int{10, 11, 23, 24}

// 	scores, err := scanForThirtyOne(hand)
// 	if len(scores) != 2 || err != nil {
// 		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q`, scores, err, 2)
// 	}
// }

// func TestNoThirtyOne(t *testing.T) {
// 	hand := []int{1, 2, 3, 4}

// 	scores, err := scanForThirtyOne(hand)
// 	if len(scores) != 0 || err != nil {
// 		t.Fatalf(`scanForThirtyOne(hand) = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestRunOfFour(t *testing.T) {
// 	hand := []int{9, 10, 11, 12}

// 	scores, err := scanForRuns(hand)
// 	if scores[0].Point != 4 || err != nil {
// 		t.Fatalf(`scanForRuns(hand) = %q, %v, want match for %#q`, scores, err, 4)
// 	}
// }

// func TestRunOfThree(t *testing.T) {
// 	hand := []int{9, 10, 11, 13}

// 	scores, err := scanForRuns(hand)
// 	if scores[0].Point != 3 || err != nil {
// 		t.Fatalf(`scanForRuns(hand) = %q, %v, want match for %#q`, scores, err, 3)
// 	}
// }

// func TestTwoRunsOfThree(t *testing.T) {
// 	hand := []int{9, 10, 11, 24}

// 	scores, err := scanForRuns(hand)
// 	if len(scores) != 2 || err != nil {
// 		t.Fatalf(`scanForRuns(hand) = %q, %v, want match for %#q`, scores, err, 2)
// 	}

// 	if scores[0].Point != 3 || err != nil {
// 		t.Fatalf(`scanForRuns(hand).points = %q, %v, want match for %#q`, scores[0].Point, err, 3)
// 	}

// 	if scores[1].Point != 3 || err != nil {
// 		t.Fatalf(`scanForRuns(hand).points = %q, %v, want match for %#q`, scores[1].Point, err, 3)
// 	}
// }

// func TestPair(t *testing.T) {
// 	hand := []int{9, 10, 11, 24}

// 	scores, err := scanForMatchingKinds(hand)
// 	utils.L(scores[0])
// 	if scores[0].Point != 2 || err != nil {
// 		t.Fatalf(`scanForMatchingKinds(hand) = %q, %v, want match for %#q`, scores, err, 2)
// 	}
// }

// func TestPairThree(t *testing.T) {
// 	hand := []int{9, 37, 11, 24}

// 	scores, err := scanForMatchingKinds(hand)
// 	if len(scores) != 3 || err != nil {
// 		t.Fatalf(`scanForMatchingKinds(hand) = %q, %v, want match for %#q`, scores, err, 3)
// 	}
// }

// func TestOneFifthteens(t *testing.T) {
// 	hand := []int{7, 8, 11, 10}

// 	scores, err := scanForFifthteens(hand)
// 	if len(scores) != 1 || err != nil {
// 		t.Fatalf(`scanForFifthteens(hand) = %q, %v, want match for %#q`, scores, err, 1)
// 	}
// }

// func TestTwoFifthteens(t *testing.T) {
// 	hand := []int{8, 9, 22, 19}

// 	scores, err := scanForFifthteens(hand)
// 	if len(scores) != 2 || err != nil {
// 		t.Fatalf(`scanForFifthteens(hand) = %q, %v, want match for %#q`, scores, err, 2)
// 	}
// }

// func TestThreeFifthteens(t *testing.T) {
// 	hand := []int{7, 8, 21, 34}

// 	scores, err := scanForFifthteens(hand)
// 	if len(scores) != 3 || err != nil {
// 		t.Fatalf(`scanForFifthteens(hand) = %q, %v, want match for %#q`, scores, err, 3)
// 	}
// }

// func TestFourFifthteens(t *testing.T) {
// 	hand := []int{7, 8, 20, 21}

// 	scores, err := scanForFifthteens(hand)
// 	if len(scores) != 4 || err != nil {
// 		t.Fatalf(`scanForFifthteens(hand) = %q, %v, want match for %#q`, scores, err, 4)
// 	}
// }

// func TestTwoFifthteensOfThree(t *testing.T) {
// 	hand := []int{7, 20, 1, 14}

// 	scores, err := scanForFifthteens(hand)
// 	if len(scores) != 2 || err != nil {
// 		t.Fatalf(`scanForFifthteens(hand) = %q, %v, want match for %#q`, scores, err, 2)
// 	}
// }

// func TestThreeFifthteensOfMix(t *testing.T) {
// 	hand := []int{7, 8, 1, 20}

// 	scores, err := scanForFifthteens(hand)
// 	if len(scores) != 3 || err != nil {
// 		t.Fatalf(`scanForFifthteens(hand) = %q, %v, want match for %#q`, scores, err, 3)
// 	}
// }

// func TestRightJack(t *testing.T) {
// 	hand := []int{7, 8, 1, 11}
// 	match := model.Match{
// 		CutGameCardId: 1,
// 	}

// 	scores, err := scanRightJackCut(hand, match)
// 	if len(scores) != 1 || err != nil {
// 		t.Fatalf(`scanRightJackCut(hand) = %q, %v, want match for %#q`, scores, err, 1)
// 	}
// }

// func TestJackOnCut(t *testing.T) {
// 	match := model.Match{
// 		CutGameCardId: 11,
// 	}

// 	scores, err := scanJackOnCut(match)
// 	if len(scores) != 1 || err != nil {
// 		t.Fatalf(`scanJackOnCut(hand) = %q, %v, want match for %#q`, scores, err, 1)
// 	}
// }

// func TestLastCard(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay:       []int{11, 24, 37, 2},
// 		CurrentPlayerTurn: 0,
// 		Players: []queries.Player{
// 			{
// 				Id:   0,
// 				Hand: []int{2, 3, 4, 5},
// 			},
// 			{
// 				Id:   1,
// 				Hand: []int{2, 3, 4, 5},
// 			},
// 		},
// 	}

// 	scores, err := scanForLastCard(match)
// 	if len(scores) != 1 || err != nil {
// 		t.Fatalf(`scanForLastCard(hand) = %q, %v, want match for %#q`, scores, err, 1)
// 	}
// }

// func TestIsNotLastCard(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay:       []int{10, 24, 37, 2},
// 		CurrentPlayerTurn: 0,
// 		Players: []queries.Player{
// 			{
// 				Id:   0,
// 				Hand: []int{2, 3, 4, 5},
// 			},
// 			{
// 				Id:   1,
// 				Hand: []int{2, 3, 4, 5},
// 			},
// 		},
// 	}

// 	scores, err := scanForLastCard(match)
// 	if len(scores) != 1 || err != nil {
// 		t.Fatalf(`scanForLastCard(hand) = %q, %v, want match for %#q`, scores, err, 1)
// 	}
// }

// func TestFlush(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay:       []int{11, 24, 37, 2},
// 		CurrentPlayerTurn: 0,
// 		Players: []queries.Player{
// 			{
// 				Id:   0,
// 				Hand: []int{2, 3, 4, 5},
// 			},
// 			{
// 				Id:   1,
// 				Hand: []int{2, 3, 4, 5},
// 			},
// 		},
// 	}

// 	scores, err := scanForFlush(match.Players[0].Hand)
// 	if len(scores) != 1 || err != nil {
// 		t.Fatalf(`scanForFlush(hand) = %q, %v, want match for %#q`, scores, err, 1)
// 	}
// }

// func TestNoFlush(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay:       []int{11, 24, 37, 2},
// 		CurrentPlayerTurn: 0,
// 		Players: []queries.Player{
// 			{
// 				Id:   0,
// 				Hand: []int{2, 3, 4, 25},
// 			},
// 			{
// 				Id:   1,
// 				Hand: []int{2, 3, 4, 5},
// 			},
// 		},
// 	}

// 	scores, err := scanForFlush(match.Players[0].Hand)
// 	if len(scores) != 0 || err != nil {
// 		t.Fatalf(`scanForFlush(hand) = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingWithFithteens(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			10,
// 			5,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{1, 2, 3, 4},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{5, 6, 7, 8},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	if len(scores.Results) != 1 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingWithRunAndFithteens(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			4,
// 			5,
// 			6,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{1, 2, 3, 4},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{5, 6, 7, 8},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	if len(scores.Results) != 2 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingLastCard(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			9,
// 			11,
// 			13,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{11, 12, 13, 24},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{25, 26, 37, 38},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	if len(scores.Results) != 1 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingLastCardAndThirtyOne(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			8,
// 			11,
// 			12,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{11, 12, 13, 24},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{25, 26, 37, 38},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	if len(scores.Results) != 2 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingMakingKinds(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			10,
// 			36,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{11, 12, 13, 24},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{25, 26, 37, 38},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	if len(scores.Results) != 1 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingNoRun(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			5,
// 			3,
// 			8,
// 			9,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{1, 2, 3, 4},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{1, 2, 3, 4},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	if len(scores.Results) != 0 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingRun(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			5,
// 			7,
// 			6,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{11, 12, 13, 24},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{25, 26, 37, 38},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	if len(scores.Results) != 1 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingRunOfFour(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay: []int{
// 			1,
// 			2,
// 			3,
// 			4,
// 		},
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{11, 12, 13, 24},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{25, 26, 37, 38},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	//Score includes a last card point
// 	if len(scores.Results) != 1 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }

// func TestPeggingJackOnCut(t *testing.T) {
// 	match := model.Match{
// 		CardsInPlay:   []int{},
// 		CutGameCardId: 11,
// 		Players: []queries.Player{
// 			{
// 				Id:   1,
// 				Hand: []int{10, 11, 12, 23},
// 			},
// 			{
// 				Id:   2,
// 				Hand: []int{25, 26, 37, 38},
// 			},
// 		},
// 	}

// 	//test pegging
// 	scores, err := countPegs(match)
// 	//Score includes a last card point
// 	if len(scores.Results) != 1 || err != nil {
// 		t.Fatalf(`countPegs() = %q, %v, want match for %#q`, scores, err, 0)
// 	}
// }
