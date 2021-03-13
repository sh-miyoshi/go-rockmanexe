package battle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

const (
	stateOpening int = iota
	stateChipSelect
	stateBeforeMain
	stateMain
	stateResult

	stateMax
)

var (
	battleCount = 0
	battleState = stateChipSelect // debug
)

// Init ...
func Init(plyr *player.Player) error {
	if err := field.Init(); err != nil {
		return fmt.Errorf("Battle field init failed: %w", err)
	}

	if err := battleplayer.Init(plyr.HP, plyr.ChipFolder); err != nil {
		return fmt.Errorf("Battle player init failed: %w", err)
	}

	return nil
}

// End ...
func End() {
	field.End()
	battleplayer.End()
}

// Process ...
func Process() {
	// TODO error handling
	anim.MgrProcess()

	switch battleState {
	case stateChipSelect:
		if battleCount == 0 {
			chipsel.Init(battleplayer.Get().ChipFolder)
			// TODO error handling
		}
		if chipsel.Process() {
			// set selected chips
			battleplayer.SetChipSelectResult(chipsel.GetSelected())
			stateChange(stateBeforeMain)
		}
	case stateBeforeMain:
		// TODO implement this
		stateChange(stateMain)
	case stateMain:
		res := battleplayer.MainProcess()
		fieldUpdates()
		if res {
			stateChange(stateChipSelect)
			return
		}
	}

	battleCount++
}

// Draw ...
func Draw() {
	if battleCount == 0 {
		// skip if initialize phase
		return
	}

	field.Draw()

	switch battleState {
	case stateChipSelect:
		battleplayer.DrawChar()
		chipsel.Draw()
	case stateMain:
		battleplayer.DrawChar()
		battleplayer.DrawChipIcon()
	}
}

func fieldUpdates() {
	p := battleplayer.Get()
	objs := []field.ObjectPosition{
		{X: p.PosX, Y: p.PosY, ID: p.ID},
	}

	// TODO set enemy pos

	field.UpdateObjectPos(objs)
}

func stateChange(nextState int) {
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next game state: %d", nextState))
	}
	battleState = nextState
	battleCount = 0
}
