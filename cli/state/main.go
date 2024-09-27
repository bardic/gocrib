package state

import (
	"github.com/bardic/cribbagev2/cli/views"
	"github.com/bardic/cribbagev2/model"
)

var ActiveMatchId int
var ActiveMatch model.Match
var ActiveViewModel views.ViewModel
var ActiveDeck model.GameDeck
var CurrentAction model.GameAction
var CurrentHandModifier model.HandModifier
