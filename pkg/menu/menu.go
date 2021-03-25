package menu

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
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

	imgBack int32
)

func Init() error {
	menuState = stateTop
	menuCount = 0

	fname := common.ImagePath + "menu/back.png"
	imgBack = dxlib.LoadGraph(fname)
	if imgBack == -1 {
		return fmt.Errorf("Failed to load menu back image %s", fname)
	}

	if err := topInit(); err != nil {
		return fmt.Errorf("Failed to init menu top: %w", err)
	}
	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	topEnd()
}

func Process() {
	switch menuState {
	case stateTop:
		topProcess()
	}
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, dxlib.TRUE)

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
