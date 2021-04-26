package menu

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
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

	if err := sound.BGMPlay(sound.BGMMenu); err != nil {
		return fmt.Errorf("failed to play bgm: %v", err)
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
	if config.Get().Debug.SkipMenu {
		return ErrGoBattle
	}

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
	return battleEnemies()
}

func stateChange(nextState int) {
	logger.Info("Change menu state from %d to %d", menuState, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	menuState = nextState
}
