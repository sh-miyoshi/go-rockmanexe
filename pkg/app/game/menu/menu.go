package menu

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateTop int = iota
	stateChipFolder
	stateGoBattle
	stateRecord
	stateNetBattle

	stateMax
)

var (
	menuState         int
	imgBack           int32
	menuFolderInst    *menuFolder
	menuRecordInst    *menuRecord
	menuNetBattleInst *menuNetBattle

	ErrGoBattle    = errors.New("go to battle")
	ErrGoNetBattle = errors.New("go to net battle")
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

	var err error
	menuFolderInst, err = folderNew(plyr)
	if err != nil {
		return fmt.Errorf("failed to init menu folder: %w", err)
	}

	if err := goBattleInit(); err != nil {
		return fmt.Errorf("failed to init menu go battle: %w", err)
	}

	menuRecordInst, err = recordNew(plyr)
	if err != nil {
		return fmt.Errorf("failed to init menu record: %w", err)
	}

	menuNetBattleInst, err = netBattleNew()
	if err != nil {
		return fmt.Errorf("failed to init menu net battle: %w", err)
	}

	if err := sound.BGMPlay(sound.BGMMenu); err != nil {
		return fmt.Errorf("failed to play bgm: %v", err)
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	topEnd()
	if menuFolderInst != nil {
		menuFolderInst.End()
	}
	goBattleEnd()
	if menuRecordInst != nil {
		menuRecordInst.End()
	}
	if menuNetBattleInst != nil {
		menuNetBattleInst.End()
	}
}

func Process() error {
	if config.Get().Debug.SkipMenu {
		return ErrGoBattle
	}

	switch menuState {
	case stateTop:
		topProcess()
	case stateChipFolder:
		menuFolderInst.Process()
	case stateGoBattle:
		if goBattleProcess() {
			return ErrGoBattle
		}
	case stateRecord:
		menuRecordInst.Process()
	case stateNetBattle:
		if menuNetBattleInst.Process() {
			return ErrGoNetBattle
		}
	}

	return nil
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, dxlib.TRUE)

	switch menuState {
	case stateTop:
		topDraw()
	case stateChipFolder:
		menuFolderInst.Draw()
	case stateGoBattle:
		goBattleDraw()
	case stateRecord:
		menuRecordInst.Draw()
	case stateNetBattle:
		menuNetBattleInst.Draw()
	}
}

func GetBattleEnemies() []enemy.EnemyParam {
	return battleEnemies()
}

func stateChange(nextState int) {
	logger.Info("Change menu state from %d to %d", menuState, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	menuState = nextState
}
