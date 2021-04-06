package menu

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
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
	imgBack   int32

	ErrGoBattle = errors.New("go to battle")
)

func Init(plyr *player.Player) error {
	menuState = stateTop

	fname := common.ImagePath + "menu/back.png"
	imgBack = dxlib.LoadGraph(fname)
	if imgBack == -1 {
		return fmt.Errorf("failed to load menu back image %s", fname)
	}

	if err := topInit(); err != nil {
		return fmt.Errorf("failed to init menu top: %w", err)
	}

	if err := folderInit(plyr); err != nil {
		return fmt.Errorf("failed to init menu folder: %w", err)
	}

	if err := goBattleInit(); err != nil {
		return fmt.Errorf("failed to init menu go battle: %w", err)
	}

	if err := recordInit(); err != nil {
		return fmt.Errorf("failed to init menu record: %w", err)
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

func Process() error {
	switch menuState {
	case stateTop:
		topProcess()
	case stateChipFolder:
		folderProcess()
	case stateGoBattle:
		if goBattleProcess() {
			return ErrGoBattle
		}
	case stateRecord:
		recordProcess()
	}

	return nil
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

func GetBattleEnemies() []enemy.EnemyParam {
	// TODO return battleEnemies()
	return []enemy.EnemyParam{
		{
			CharID: enemy.IDTarget,
			PosX:   4,
			PosY:   1,
			HP:     1000,
		},
		{
			CharID: enemy.IDTarget,
			PosX:   5,
			PosY:   1,
			HP:     1000,
		},
	}
}

func stateChange(nextState int) {
	logger.Info("Change menu state from %d to %d", menuState, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	menuState = nextState
}
