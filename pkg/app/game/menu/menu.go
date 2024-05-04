package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	stateTop int = iota
	stateChipFolder
	stateGoBattle
	statePlayerStatus
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
	ResultGoTalkAI
)

type menuStateInstance interface {
	End()
	Process() bool
	Draw()
	GetResult() Result // TODO: WIP
}

var (
	menuState     int
	imgBack       int
	currentInst   menuStateInstance
	playerInfo    *player.Player
	battleEnemies []enemy.EnemyParam
)

func Init(plyr *player.Player) error {
	menuState = stateTop
	playerInfo = plyr
	currentInst = nil

	fname := config.ImagePath + "menu/back.png"
	if imgBack = dxlib.LoadGraph(fname); imgBack == -1 {
		return fmt.Errorf("failed to load menu back image %s", fname)
	}

	if err := sound.BGMPlay(sound.BGMMenu); err != nil {
		return fmt.Errorf("failed to play bgm: %v", err)
	}

	if config.Get().Debug.SkipMenu {
		// Start from battle mode for debug, so set debug data
		battleEnemies = []enemy.EnemyParam{
			{
				CharID: enemy.IDTarget,
				Pos:    point.Point{X: 4, Y: 1},
				HP:     1000,
			},
		}
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	if currentInst != nil {
		currentInst.End()
		currentInst = nil
	}
	// TODO goBattleEnd()
}

func Process() (Result, error) {
	if config.Get().Debug.SkipMenu {
		return ResultGoBattle, nil
	}

	var err error
	switch menuState {
	case stateTop:
		if currentInst == nil {
			currentInst, err = topNew(playerInfo)
			if err != nil {
				return ResultContinue, err
			}
		}

		if currentInst.Process() {
			if res := currentInst.GetResult(); res != ResultContinue {
				return res, nil
			}
			next := currentInst.(*menuTop).GetNextState()
			stateChange(next)
		}
	case stateChipFolder:
		if currentInst == nil {
			currentInst, err = folderNew(playerInfo)
			if err != nil {
				return ResultContinue, err
			}
		}

		if currentInst.Process() {
			stateChange(stateTop)
		}
	case stateGoBattle:
		if currentInst == nil {
			currentInst, err = goBattleNew()
			if err != nil {
				return ResultContinue, err
			}
		}

		if currentInst.Process() {
			if res := currentInst.GetResult(); res != ResultContinue {
				return res, nil
			}
			stateChange(stateTop)
		}
	case statePlayerStatus:
		if currentInst == nil {
			currentInst, err = playerStatusNew(playerInfo)
			if err != nil {
				return ResultContinue, err
			}
		}

		if currentInst.Process() {
			stateChange(stateTop)
		}
	case stateNetBattle:
		if currentInst == nil {
			currentInst, err = netBattleNew()
			if err != nil {
				return ResultContinue, err
			}
		}

		if currentInst.Process() {
			if res := currentInst.GetResult(); res != ResultContinue {
				return res, nil
			}
			stateChange(stateTop)
		}
	case stateInvalidChip:
		if currentInst == nil {
			currentInst, err = invalidChipNew(playerInfo)
			if err != nil {
				return ResultContinue, err
			}
		}

		if currentInst.Process() {
			stateChange(stateTop)
		}
	case stateDevFeature:
		if currentInst == nil {
			currentInst, err = devFeatureNew()
			if err != nil {
				return ResultContinue, err
			}
		}

		if currentInst.Process() {
			return currentInst.GetResult(), nil
		}
	}

	return ResultContinue, nil
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, true)
	if currentInst != nil {
		currentInst.Draw()
	}
}

func GetBattleEnemies() []enemy.EnemyParam {
	return battleEnemies
}

func stateChange(nextState int) {
	logger.Info("Change menu state from %d to %d", menuState, nextState)
	if nextState < 0 || nextState >= stateMax {
		system.SetError(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	menuState = nextState
	if currentInst != nil {
		currentInst.End()
		currentInst = nil
	}
}
