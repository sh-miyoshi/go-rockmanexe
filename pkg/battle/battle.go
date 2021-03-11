package battle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
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
	battleState = stateMain // debug
)

// Init ...
func Init(plyr *player.Player) error {
	if err := field.Init(); err != nil {
		return fmt.Errorf("Battle field init failed: %w", err)
	}

	if err := battleplayer.Init(plyr.HP); err != nil {
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
	case stateMain:
		battleplayer.MainProcess()
		fieldUpdates()
	}
}

// Draw ...
func Draw() {
	field.Draw()
	battleplayer.Draw()
}

func fieldUpdates() {
	px, py := battleplayer.GetPos()
	objs := []field.ObjectPosition{
		{X: px, Y: py, ID: battleplayer.GetID()},
	}

	// TODO set enemy pos

	field.UpdateObjectPos(objs)
}
