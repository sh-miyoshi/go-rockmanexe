package battle

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/opening"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/win"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

const (
	stateOpening int = iota
	stateChipSelect
	stateBeforeMain
	stateMain
	stateResultWin
	stateResultLose

	stateMax
)

var (
	battleCount int
	battleState int
	playerInst  *battleplayer.BattlePlayer
	enemyList   []enemy.EnemyParam
	gameCount   int
	b4mainInst  *titlemsg.TitleMsg
	loseInst    *titlemsg.TitleMsg

	ErrWin  = errors.New("player win")
	ErrLose = errors.New("playser lose")
)

// Init ...
func Init(plyr *player.Player, enemies []enemy.EnemyParam) error {
	enemyList = enemies
	gameCount = 0
	battleCount = 0
	battleState = stateOpening
	b4mainInst = nil
	loseInst = nil

	var err error
	playerInst, err = battleplayer.New(plyr)
	if err != nil {
		return fmt.Errorf("battle player init failed: %w", err)
	}
	anim.New(playerInst)

	if err := field.Init(); err != nil {
		return fmt.Errorf("battle field init failed: %w", err)
	}

	if err := skill.Init(); err != nil {
		return fmt.Errorf("skill init failed: %w", err)
	}

	if err := effect.Init(); err != nil {
		return fmt.Errorf("effect init failed: %w", err)
	}

	return nil
}

// End ...
func End() {
	anim.Cleanup()
	damage.RemoveAll()
	field.End()
	playerInst.End()
	skill.End()
	enemy.End()
	effect.End()
	win.End()
}

// Process ...
func Process() error {
	switch battleState {
	case stateOpening:
		if battleCount == 0 {
			if err := opening.Init(enemyList); err != nil {
				return fmt.Errorf("opening init failed: %w", err)
			}
		}

		if opening.Process() {
			opening.End()
			if err := enemy.Init(playerInst.ID, enemyList); err != nil {
				return fmt.Errorf("enemy init failed: %w", err)
			}
			stateChange(stateChipSelect)
			return nil
		}
	case stateChipSelect:
		if battleCount == 0 {
			if err := chipsel.Init(playerInst.ChipFolder); err != nil {
				return fmt.Errorf("failed to initialize chip select: %w", err)
			}
		}
		if chipsel.Process() {
			// set selected chips
			playerInst.SetChipSelectResult(chipsel.GetSelected())
			stateChange(stateBeforeMain)
			return nil
		}
	case stateBeforeMain:
		if battleCount == 0 {
			fname := common.ImagePath + "battle/msg_start.png"
			var err error
			b4mainInst, err = titlemsg.New(fname)
			if err != nil {
				return fmt.Errorf("failed to initialize before main: %w", err)
			}
		}

		if b4mainInst.Process() {
			b4mainInst.End()
			stateChange(stateMain)
			return nil
		}
	case stateMain:
		gameCount++
		if err := anim.MgrProcess(true); err != nil {
			return fmt.Errorf("failed to handle animation: %w", err)
		}

		switch playerInst.NextAction {
		case battleplayer.NextActChipSelect:
			stateChange(stateChipSelect)
			playerInst.NextAction = battleplayer.NextActNone
			return nil
		case battleplayer.NextActLose:
			stateChange(stateResultLose)
			return nil
		}
		if err := enemy.MgrProcess(); err != nil {
			if errors.Is(err, enemy.ErrGameEnd) {
				playerInst.EnableAct = false
				stateChange(stateResultWin)
				return nil
			}
			return fmt.Errorf("failed to process enemy: %w", err)
		}
		fieldUpdates()
	case stateResultWin:
		if battleCount == 0 {
			if err := win.Init(gameCount); err != nil {
				return fmt.Errorf("failed to initialize result win: %w", err)
			}
		}

		if err := anim.MgrProcess(false); err != nil {
			return fmt.Errorf("failed to handle animation: %w", err)
		}

		if win.Process() {
			return ErrWin
		}
	case stateResultLose:
		if battleCount == 0 {
			fname := common.ImagePath + "battle/msg_lose.png"
			var err error
			loseInst, err = titlemsg.New(fname)
			if err != nil {
				return fmt.Errorf("failed to initialize lose: %w", err)
			}
		}

		if loseInst.Process() {
			loseInst.End()
			return ErrLose
		}
	}

	battleCount++
	return nil
}

// Draw ...
func Draw() {
	field.Draw()
	anim.MgrDraw()

	switch battleState {
	case stateOpening:
		opening.Draw()
	case stateChipSelect:
		playerInst.DrawFrame(true, false)
		chipsel.Draw()
	case stateBeforeMain:
		playerInst.DrawFrame(false, true)
		if b4mainInst != nil {
			b4mainInst.Draw()
		}
	case stateMain:
		playerInst.DrawFrame(false, true)
	case stateResultWin:
		playerInst.DrawFrame(false, true)
		win.Draw()
	case stateResultLose:
		playerInst.DrawFrame(false, false)
		if loseInst != nil {
			loseInst.Draw()
		}
	}
}

func fieldUpdates() {
	objs := []field.ObjectPosition{
		{X: playerInst.PosX, Y: playerInst.PosY, ID: playerInst.ID},
	}

	enemies := enemy.GetEnemyPositions()
	objs = append(objs, enemies...)

	field.UpdateObjectPos(objs)
}

func stateChange(nextState int) {
	logger.Info("Change battle state from %d to %d", battleState, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	battleState = nextState
	battleCount = 0
}
