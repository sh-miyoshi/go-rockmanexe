package battle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
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
	// TODO error handling
	anim.MgrProcess()

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

func moveObject(x, y *int, direct int, isMove bool) bool {
	nx := *x
	ny := *y

	switch direct {
	case common.DirectUp:
		if ny <= 0 {
			return false
		}
		ny--
	case common.DirectDown:
		if ny >= fieldNumY-1 {
			return false
		}
		ny++
	case common.DirectLeft:
		if nx <= 0 {
			return false
		}
		nx--
	case common.DirectRight:
		if nx >= fieldNumX-1 {
			return false
		}
		nx++
	}

	// TODO field panel is player?

	if isMove {
		*x = nx
		*y = ny
	}

	return true
}
