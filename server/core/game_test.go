package core

import (
	"testing"

	"github.com/bardic/gocrib/queries/queries"
	"github.com/bardic/gocrib/vo"
)

func TestThirtyOne(t *testing.T) {

	expectedPointSets := 1
	expectedPoints := 2

	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueAce,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueQueen,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scores, err := scanForThirtyOne(hand)
	if len(scores) != expectedPointSets || scores[0].Point == &expectedPoints || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %v, %v, want match for %#q`, scores, err, 2)
	}
}

func TestNoThirtyOne(t *testing.T) {
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueAce,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scores, err := scanForThirtyOne(hand)
	if len(scores) != 0 || err != nil {
		t.Fatalf(`scanForThirtyOne(hand) = %v, %v, want match for %#q`, scores, err, 0)
	}
}

func TestRunOfFour(t *testing.T) {

	expectedPoints := 4

	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueQueen,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scores, err := scanForRuns(hand)
	if *scores[0].Point != expectedPoints || err != nil {
		t.Fatalf(`scanForRuns(hand) = %v, %v, want match for %v`, scores[0].Point, err, expectedPoints)
	}
}

func TestRunOfThree(t *testing.T) {

	expectedPoints := 3
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueKing,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scores, err := scanForRuns(hand)
	if *scores[0].Point != expectedPoints || err != nil {
		t.Fatalf(`scanForRuns(hand) = %v, %v, want match for %#q`, scores, err, 3)
	}
}

func TestTwoRunsOfThree(t *testing.T) {
	expectedPoints := 3
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitClubs,
			},
		},
	}

	scores, err := scanForRuns(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForRuns(hand) = %v, %v, want match for %#q`, scores, err, 2)
	}

	if *scores[0].Point != expectedPoints || err != nil {
		t.Fatalf(`scanForRuns(hand).points = %v, %v, want match for %#q`, scores[0].Point, err, 3)
	}

	if *scores[1].Point != expectedPoints || err != nil {
		t.Fatalf(`scanForRuns(hand).points = %v, %v, want match for %#q`, scores[1].Point, err, 3)
	}
}

func TestPair(t *testing.T) {
	expectedPoints := 2
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitClubs,
			},
		},
	}

	scores, err := scanForMatchingKinds(hand)
	if *scores[0].Point != expectedPoints || err != nil {
		t.Fatalf(`scanForMatchingKinds(hand) = %v, %v, want match for %#q`, scores, err, 2)
	}
}

func TestTwoPair(t *testing.T) {
	// expectedPoints := 3
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitClubs,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitClubs,
			},
		},
	}

	scores, err := scanForMatchingKinds(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForMatchingKinds(hand) = %v, %v, want match for %#q`, scores, err, 3)
	}
}

func TestOneFifteens(t *testing.T) {
	// expectedPoints := 3
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitClubs,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitClubs,
			},
		},
	}

	scores, err := scanForFifteens(hand)
	if len(scores) != 1 || err != nil {
		t.Fatalf(`scanForFifteens(hand) = %v, %v, want match for %#q`, scores, err, 1)
	}
}

func TestTwoFifteens(t *testing.T) {
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitClubs,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitClubs,
			},
		},
	}

	scores, err := scanForFifteens(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForFifteens(hand) = %v, %v, want match for %#q`, scores, err, 2)
	}
}

func TestThreeFifteens(t *testing.T) {
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitClubs,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitSpades,
			},
		},
	}

	scores, err := scanForFifteens(hand)
	if len(scores) != 3 || err != nil {
		t.Fatalf(`scanForFifteens(hand) = %v, %v, want match for %#q`, scores, err, 3)
	}
}

func TestFourFifteens(t *testing.T) {
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueFive,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueFive,
				Suit:  queries.CardsuitClubs,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueFive,
				Suit:  queries.CardsuitDiamonds,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueFive,
				Suit:  queries.CardsuitSpades,
			},
		},
	}

	scores, err := scanForFifteens(hand)
	if len(scores) != 4 || err != nil {
		t.Fatalf(`scanForFifteens(hand) = %v, %v, want match for %#q`, scores, err, 4)
	}
}

func TestTwoFifteensOfThree(t *testing.T) {
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueFour,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitClubs,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueFive,
				Suit:  queries.CardsuitSpades,
			},
		},
	}

	scores, err := scanForFifteens(hand)
	if len(scores) != 2 || err != nil {
		t.Fatalf(`scanForFifteens(hand) = %v, %v, want match for %#q`, scores, err, 2)
	}
}

func TestRightJack(t *testing.T) {
	hand := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitClubs,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitDiamonds,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitSpades,
			},
		},
	}

	cutCard := vo.GameCard{
		Card: queries.Card{
			Value: queries.CardvalueAce,
			Suit:  queries.CardsuitSpades,
		},
	}

	scores, err := scanRightJackCut(hand, cutCard)
	if len(scores) != 1 || err != nil {
		t.Fatalf(`scanRightJackCut(hand) = %v, %v, want match for %#q`, scores, err, 1)
	}
}

func TestJackOnCut(t *testing.T) {
	cutCard := vo.GameCard{
		Card: queries.Card{
			Value: queries.CardvalueJack,
			Suit:  queries.CardsuitSpades,
		},
	}

	scores, err := scanJackOnCut(cutCard)
	if len(scores) != 1 || err != nil {
		t.Fatalf(`scanJackOnCut(hand) = %v, %v, want match for %#q`, scores, err, 1)
	}
}

func TestLastCard(t *testing.T) {
	// cardsInPlay := []vo.GameCard{
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueNine,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueTen,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueTen,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// }

	// p1Id := 1
	// p2Id := 2

	// players := []vo.GamePlayer{
	// 	{
	// 		Player: queries.Player{
	// 			ID: &p1Id,
	// 		},
	// 		Hand: []vo.GameCard{
	// 			{
	// 				Card: queries.Card{
	// 					Value: queries.CardvalueTen,
	// 					Suit:  queries.CardsuitHearts,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Player: queries.Player{
	// 			ID: &p2Id,
	// 		},
	// 		Hand: []vo.GameCard{
	// 			{
	// 				Card: queries.Card{
	// 					Value: queries.CardvalueNine,
	// 					Suit:  queries.CardsuitHearts,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// scores, err := scanForLastCard(p1Id, cardsInPlay, players)

	// if err != nil {
	// 	t.Fatalf(`scanForLastCard(hand) = %v, %v, want match for %#q`, scores, err, 1)
	// }

	// if len(scores) > 0 && *scores[0].Point == 2 {
	// 	t.Fatalf(`scanForLastCard(hand) = %v, %v, want match for %#q`, scores, err, 1)
	// }

	// if len(scores) > 0 && *scores[0].Point == 0 {
	// 	t.Fatalf(`scanForLastCard(hand) = %v, %v, want match for %#q`, scores, err, 1)
	// }
}

func TestFlush(t *testing.T) {
	cardsInPlay := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueNine,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
	}
	scores, err := scanForFlush(cardsInPlay)
	if len(scores) != 1 || err != nil {
		t.Fatalf(`scanForFlush(hand) = %v, %v, want match for %#q`, scores, err, 1)
	}
}

func TestPeggingWithFifteens(t *testing.T) {
	cardsInPlay := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueEight,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSeven,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scoreMatch := vo.ScoreMatch{
		ActivePlayerId: nil,
		PlayerSeekId:   nil,
		CardsInPlay:    &cardsInPlay,
		Players:        nil,
	}

	//test pegging
	scores, err := countPegs(scoreMatch)
	if len(scores.Results) != 1 || err != nil {
		t.Fatalf(`countPegs() = %v, %v, want match for %#q`, scores, err, 0)
	}
}

func TestPeggingWithRunAndFifteens(t *testing.T) {
	cardsInPlay := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueFour,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueFive,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueSix,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scoreMatch := vo.ScoreMatch{
		ActivePlayerId: nil,
		PlayerSeekId:   nil,
		CardsInPlay:    &cardsInPlay,
		Players:        nil,
	}

	//test pegging
	scores, err := countPegs(scoreMatch)
	if len(scores.Results) != 2 || err != nil {
		t.Fatalf(`countPegs() = %v, %v, want match for %#q`, scores, err, 0)
	}
}

func TestPeggingLastCard(t *testing.T) {
	// cardsInPlay := []vo.GameCard{
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueNine,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueTen,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueKing,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// }

	// p1Id := 1
	// p2Id := 2

	// players := []vo.GamePlayer{
	// 	{
	// 		Player: queries.Player{
	// 			ID: &p1Id,
	// 		},
	// 		Hand: []vo.GameCard{
	// 			{
	// 				Card: queries.Card{
	// 					Value: queries.CardvalueTen,
	// 					Suit:  queries.CardsuitHearts,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Player: queries.Player{
	// 			ID: &p2Id,
	// 		},
	// 		Hand: []vo.GameCard{
	// 			{
	// 				Card: queries.Card{
	// 					Value: queries.CardvalueNine,
	// 					Suit:  queries.CardsuitHearts,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// scoreMatch := vo.ScoreMatch{
	// 	ActivePlayerId: &p1Id,
	// 	PlayerSeekId:   &p1Id,
	// 	CardsInPlay:    &cardsInPlay,
	// 	Players:        &players,
	// }

	// //test pegging
	// scores, err := countPegs(scoreMatch)
	// if len(scores.Results) != 1 || err != nil {
	// 	t.Fatalf(`countPegs() = %v, %v, want match for %#q`, scores, err, 0)
	// }
}

func TestPeggingLastCardAndThirtyOne(t *testing.T) {
	// cardsInPlay := []vo.GameCard{
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueTen,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueQueen,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueKing,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// 	{
	// 		Card: queries.Card{
	// 			Value: queries.CardvalueAce,
	// 			Suit:  queries.CardsuitHearts,
	// 		},
	// 	},
	// }

	// p1Id := 1
	// p2Id := 2

	// players := []vo.GamePlayer{
	// 	{
	// 		Player: queries.Player{
	// 			ID: &p1Id,
	// 		},
	// 		Hand: []vo.GameCard{
	// 			{
	// 				Card: queries.Card{
	// 					Value: queries.CardvalueTen,
	// 					Suit:  queries.CardsuitHearts,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	{
	// 		Player: queries.Player{
	// 			ID: &p2Id,
	// 		},
	// 		Hand: []vo.GameCard{
	// 			{
	// 				Card: queries.Card{
	// 					Value: queries.CardvalueNine,
	// 					Suit:  queries.CardsuitHearts,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// scoreMatch := vo.ScoreMatch{
	// 	ActivePlayerId: &p1Id,
	// 	PlayerSeekId:   &p1Id,
	// 	CardsInPlay:    &cardsInPlay,
	// 	Players:        &players,
	// }

	// //test pegging
	// scores, err := countPegs(scoreMatch)
	// if len(scores.Results) != 2 || err != nil {
	// 	t.Fatalf(`countPegs() = %v, %v, want match for %#q`, scores, err, 0)
	// }
}

func TestPeggingMakingKinds(t *testing.T) {
	cardsInPlay := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueKing,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scoreMatch := vo.ScoreMatch{
		ActivePlayerId: nil,
		PlayerSeekId:   nil,
		CardsInPlay:    &cardsInPlay,
		Players:        nil,
	}

	//test pegging
	scores, err := countPegs(scoreMatch)
	if len(scores.Results) != 1 || err != nil {
		t.Fatalf(`countPegs() = %v, %v, want match for %#q`, scores, err, 0)
	}
}

func TestPeggingRun(t *testing.T) {
	cardsInPlay := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueQueen,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scoreMatch := vo.ScoreMatch{
		ActivePlayerId: nil,
		PlayerSeekId:   nil,
		CardsInPlay:    &cardsInPlay,
		Players:        nil,
	}

	//test pegging
	scores, err := countPegs(scoreMatch)
	if len(scores.Results) != 1 || err != nil {
		t.Fatalf(`countPegs() = %v, %v, want match for %#q`, scores, err, 0)
	}
}

func TestPeggingRunOfFour(t *testing.T) {
	cardsInPlay := []vo.GameCard{
		{
			Card: queries.Card{
				Value: queries.CardvalueTen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueJack,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueQueen,
				Suit:  queries.CardsuitHearts,
			},
		},
		{
			Card: queries.Card{
				Value: queries.CardvalueKing,
				Suit:  queries.CardsuitHearts,
			},
		},
	}

	scoreMatch := vo.ScoreMatch{
		ActivePlayerId: nil,
		PlayerSeekId:   nil,
		CardsInPlay:    &cardsInPlay,
		Players:        nil,
	}

	//test pegging
	scores, err := countPegs(scoreMatch)
	if len(scores.Results) != 1 || err != nil {
		t.Fatalf(`countPegs() = %v, %v, want match for %#q`, scores, err, 0)
	}
}
