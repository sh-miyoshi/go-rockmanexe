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
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/skill"
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
	battleState = stateChipSelect // debug
	playerInst  *battleplayer.BattlePlayer
)

// Init ...
func Init(plyr *player.Player) error {
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

	if err := enemy.Init(playerInst.ID); err != nil {
		return fmt.Errorf("Enemy init failed: %w", err)
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
		}
	case stateBeforeMain:
		// TODO implement this
		stateChange(stateMain)
	case stateMain:
		if err := anim.MgrProcess(); err != nil {
			return fmt.Errorf("Failed to handle animation: %w", err)
		}

		switch playerInst.NextAction {
		case battleplayer.NextActChipSelect:
			stateChange(stateChipSelect)
		case battleplayer.MextActLose:
			// TODO return lose
		}
		if err := enemy.MgrProcess(); err != nil {
			if errors.Is(err, enemy.ErrGameEnd) {
				// TODO return win
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
	if battleCount == 0 {
		// skip if initialize phase
		return
	}

	field.Draw()
	anim.MgrDraw()

	switch battleState {
	case stateChipSelect:
		playerInst.DrawFrame(true, false)
		enemy.MgrDraw()
		chipsel.Draw()
	case stateMain:
		playerInst.DrawFrame(false, true)
		enemy.MgrDraw()
	}
}

func fieldUpdates() {
	objs := []field.ObjectPosition{
		{X: playerInst.PosX, Y: playerInst.PosY, ID: playerInst.ID},
	}

	enemies := enemy.GetEnemies()
	for _, e := range enemies {
		objs = append(objs, field.ObjectPosition{
			X:  e.PosX,
			Y:  e.PosY,
			ID: e.ID,
		})
	}

	field.UpdateObjectPos(objs)
}

func stateChange(nextState int) {
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next game state: %d", nextState))
	}
	battleState = nextState
	battleCount = 0
}
