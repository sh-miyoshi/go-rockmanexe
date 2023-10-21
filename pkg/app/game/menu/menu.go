package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
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
	stateDevFeature

	stateMax
)

type Result int

const (
	ResultContinue Result = iota
	ResultGoBattle
	ResultGoNetBattle
	ResultGoMap
	ResultGoScratch
	ResultGoNaviCustom
)

var (
	menuState           int
	imgBack             int
	menuTopInst         *menuTop
	menuFolderInst      *menuFolder
	menuRecordInst      *menuRecord
	menuNetBattleInst   *menuNetBattle
	menuInvalidChipInst *menuInvalidChip
	menuDevFeatureInst  *menuDevFeature
	specificEnemy       []enemy.EnemyParam
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

	menuDevFeatureInst, err = devFeatureNew()
	if err != nil {
		return fmt.Errorf("failed to init menu dev feature: %w", err)
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
	if menuDevFeatureInst != nil {
		menuDevFeatureInst.End()
	}
}

func Process() (Result, error) {
	if config.Get().Debug.SkipMenu {
		return ResultGoMap, nil
	}

	switch menuState {
	case stateTop:
		res := menuTopInst.Process()
		if res != ResultContinue {
			return res, nil
		}
	case stateChipFolder:
		menuFolderInst.Process()
	case stateGoBattle:
		if goBattleProcess() {
			return ResultGoBattle, nil
		}
	case stateRecord:
		menuRecordInst.Process()
	case stateNetBattle:
		if menuNetBattleInst.Process() {
			return ResultGoNetBattle, nil
		}
	case stateInvalidChip:
		menuInvalidChipInst.Process()
	case stateDevFeature:
		return menuDevFeatureInst.Process()
	}

	return ResultContinue, nil
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, true)

	switch menuState {
	case stateTop:
		menuTopInst.Draw()
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
	case stateDevFeature:
		menuDevFeatureInst.Draw()
	}
}

func GetBattleEnemies() []enemy.EnemyParam {
	if len(specificEnemy) > 0 {
		return specificEnemy
	}
	return battleEnemies()
}

func stateChange(nextState int) {
	logger.Info("Change menu state from %d to %d", menuState, nextState)
	if nextState < 0 || nextState >= stateMax {
		common.SetError(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	menuState = nextState
}
