package battle

import (
	"fmt"

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
	if err := fieldInit(); err != nil {
		return fmt.Errorf("Battle field init failed: %w", err)
	}

	if err := playerInit(plyr.HP); err != nil {
		return fmt.Errorf("Battle player init failed: %w", err)
	}

	return nil
}

// End ...
func End() {
	fieldEnd()
	playerEnd()
}

// Process ...
func Process() {
	switch battleState {
	case stateMain:
		playerMainProcess()
	}
}

// Draw ...
func Draw() {
	fieldDraw()
	playerDraw()
}
