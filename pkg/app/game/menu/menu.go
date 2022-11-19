package menu

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateTop int = iota
	stateChipFolder
	stateGoBattle
	stateRecord
	stateNetBattle
	stateInvalidChip

	stateMax
)

var (
	menuState           int
	imgBack             int
	menuTopInst         *menuTop
	menuFolderInst      *menuFolder
	menuRecordInst      *menuRecord
	menuNetBattleInst   *menuNetBattle
	menuInvalidChipInst *menuInvalidChip

	ErrGoBattle    = errors.New("go to battle")
	ErrGoNetBattle = errors.New("go to net battle")
	ErrGoMap       = errors.New("go to map")
)

func Init(plyr *player.Player) error {
	menuState = stateTop

	fname := common.ImagePath + "menu/back.png"
	imgBack = dxlib.LoadGraph(fname)
	if imgBack == -1 {
		return fmt.Errorf("failed to load menu back image %s", fname)
	}

	var err error
	menuTopInst, err = topNew(plyr)
	if err != nil {
		return fmt.Errorf("failed to init menu top: %w", err)
	}

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

	menuInvalidChipInst, err = invalidChipNew()
	if err != nil {
		return fmt.Errorf("failed to init menu invalid chip: %w", err)
	}

	if err := sound.BGMPlay(sound.BGMMenu); err != nil {
		return fmt.Errorf("failed to play bgm: %v", err)
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	if menuTopInst != nil {
		menuTopInst.End()
	}
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
	if menuInvalidChipInst != nil {
		menuInvalidChipInst.End()
	}
}

func Process() error {
	if config.Get().Debug.SkipMenu {
		return ErrGoBattle
	}

	switch menuState {
	case stateTop:
		menuTopInst.Process()

		if config.Get().Debug.EnableDevFeature {
			if inputs.CheckKey(inputs.KeyLButton) == 1 {
				return ErrGoMap
			}
			if inputs.CheckKey(inputs.KeyRButton) == 1 {
				field.Set4x4Area()
				return ErrGoBattle
			}
		}
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
	case stateInvalidChip:
		menuInvalidChipInst.Process()
	}

	return nil
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, true)

	switch menuState {
	case stateTop:
		menuTopInst.Draw()

		if config.Get().Debug.EnableDevFeature {
			draw.String(50, 220, 0x000000, "Debug機能")
			draw.String(65, 250, 0x000000, "L-btn: マップ移動")
			draw.String(65, 275, 0x000000, "R-btn: 4x4 対戦")
		}
	case stateChipFolder:
		menuFolderInst.Draw()
	case stateGoBattle:
		goBattleDraw()
	case stateRecord:
		menuRecordInst.Draw()
	case stateNetBattle:
		menuNetBattleInst.Draw()
	case stateInvalidChip:
		menuInvalidChipInst.Draw()
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
