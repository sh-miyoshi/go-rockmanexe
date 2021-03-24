package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateTop int = iota
	stateChipFolder
	stateGoBattle
	stateRecord

	stateMax
)

var (
	menuState int
	menuCount int
)

func Init() {
	menuState = stateTop
	menuCount = 0

	topInit()
}

func End() {
	topEnd()
}

func Process() {
	switch menuState {
	case stateTop:
		topProcess()
	}
}

func Draw() {
	switch menuState {
	case stateTop:
		topDraw()
	}
}

func stateChange(nextState int) {
	logger.Info("Change menu state from %d to %d", menuState, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	menuState = nextState
	menuCount = 0
}
