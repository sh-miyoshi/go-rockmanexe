package battle

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/b4main"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/opening"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/win"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

const (
	stateOpening int = iota
	stateChipSelect
	stateBeforeMain
	stateMain
	stateResultWin

	stateMax
)

var (
	battleCount int
	battleState int
	playerInst  *battleplayer.BattlePlayer
	enemyList   []enemy.EnemyParam
	gameCount   int

	ErrWin  = errors.New("player win")
	ErrLose = errors.New("playser lose")
)

// Init ...
func Init(plyr *player.Player, enemies []enemy.EnemyParam) error {
	enemyList = enemies
	gameCount = 0
	battleCount = 0
	battleState = stateOpening

	var err error
	playerInst, err = battleplayer.New(plyr)
	if err != nil {
		return fmt.Errorf("Battle player init failed: %w", err)
	}
	anim.New(playerInst)

	if err := field.Init(); err != nil {
		return fmt.Errorf("Battle field init failed: %w", err)
	}

	if err := skill.Init(); err != nil {
		return fmt.Errorf("Skill init failed: %w", err)
	}

	if err := effect.Init(); err != nil {
		return fmt.Errorf("Effect init failed: %w", err)
	}

	return nil
}

// End ...
func End() {
	anim.Cleanup()
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
				return fmt.Errorf("Opening init failed: %w", err)
			}
		}

		if opening.Process() {
			opening.End()
			if err := enemy.Init(playerInst.ID, enemyList); err != nil {
				return fmt.Errorf("Enemy init failed: %w", err)
			}
			stateChange(stateChipSelect) // debug
			return nil
		}
	case stateChipSelect:
		if battleCount == 0 {
			if err := chipsel.Init(playerInst.ChipFolder); err != nil {
				return fmt.Errorf("Failed to initialize chip select: %w", err)
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
			if err := b4main.Init(); err != nil {
				return fmt.Errorf("Failed to initialize before main: %w", err)
			}
		}

		if b4main.Process() {
			b4main.End()
			stateChange(stateMain)
			return nil
		}
	case stateMain:
		gameCount++
		if err := anim.MgrProcess(true); err != nil {
			return fmt.Errorf("Failed to handle animation: %w", err)
		}

		switch playerInst.NextAction {
		case battleplayer.NextActChipSelect:
			stateChange(stateChipSelect)
			playerInst.NextAction = battleplayer.NextActNone
			return nil
		case battleplayer.MextActLose:
			return ErrLose
		}
		if err := enemy.MgrProcess(); err != nil {
			if errors.Is(err, enemy.ErrGameEnd) {
				playerInst.EnableAct = false
				stateChange(stateResultWin)
				return nil
			}
			return fmt.Errorf("Failed to process enemy: %w", err)
		}
		fieldUpdates()
	case stateResultWin:
		if battleCount == 0 {
			if err := win.Init(gameCount); err != nil {
				return fmt.Errorf("Failed to initialize result win: %w", err)
			}
		}

		if err := anim.MgrProcess(false); err != nil {
			return fmt.Errorf("Failed to handle animation: %w", err)
		}

		if win.Process() {
			return ErrWin
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
		b4main.Draw()
	case stateMain:
		playerInst.DrawFrame(false, true)
	case stateResultWin:
		playerInst.DrawFrame(false, true)
		win.Draw()
	}
}

func fieldUpdates() {
	objs := []field.ObjectPosition{
		{X: playerInst.PosX, Y: playerInst.PosY, ID: playerInst.ID},
	}

	enemies := enemy.GetEnemyPositions()
	for _, e := range enemies {
		objs = append(objs, e)
	}

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
