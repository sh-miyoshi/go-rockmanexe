package battle

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/b4main"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/opening"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

const (
	stateOpening int = iota
	stateChipSelect
	stateBeforeMain
	stateMain
	stateResult

	stateMax
)

var (
	battleCount = 0
	battleState = stateOpening // debug
	playerInst  *battleplayer.BattlePlayer
	enemyList   []enemy.EnemyParam

	ErrWin  = errors.New("player win")
	ErrLose = errors.New("playser lose")
)

// Init ...
func Init(plyr *player.Player, enemies []enemy.EnemyParam) error {
	enemyList = enemies

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
	field.End()
	playerInst.End()
	skill.End()
	enemy.End()
	effect.End()
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
			stateChange(stateChipSelect)
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
		if err := anim.MgrProcess(); err != nil {
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
				return ErrWin
			}
			return fmt.Errorf("Failed to process enemy: %w", err)
		}
		fieldUpdates()

		damage.MgrProcess()
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
