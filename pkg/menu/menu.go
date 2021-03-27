package menu

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
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

func Init(plyr *player.Player) error {
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

	if err := folderInit(plyr); err != nil {
		return fmt.Errorf("Failed to init menu folder: %w", err)
	}

	if err := goBattleInit(); err != nil {
		return fmt.Errorf("Failed to init menu go battle: %w", err)
	}

	if err := recordInit(); err != nil {
		return fmt.Errorf("Failed to init menu record: %w", err)
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	topEnd()
	folderEnd()
	goBattleEnd()
	recordEnd()
}

func Process() {
	switch menuState {
	case stateTop:
		topProcess()
	case stateChipFolder:
		folderProcess()
	case stateGoBattle:
		goBattleProcess()
	case stateRecord:
		recordProcess()
	}
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, dxlib.TRUE)

	switch menuState {
	case stateTop:
		topDraw()
	case stateChipFolder:
		folderDraw()
	case stateGoBattle:
		goBattleDraw()
	case stateRecord:
		recordDraw()
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
