package state

import (
	"github.com/bardic/cribbagev2/model"
)

var PlayerId int
var ActiveMatchId int
var ActiveMatch model.GameMatch
var ActiveDeck model.GameDeck
var CurrentAction model.GameAction
var CurrentHandModifier model.HandModifier
