package battle

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/chipsel"
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
)

// Init ...
func Init(plyr *player.Player) error {
	if err := field.Init(); err != nil {
		return fmt.Errorf("Battle field init failed: %w", err)
	}

	if err := battleplayer.Init(plyr.HP, plyr.ChipFolder); err != nil {
		return fmt.Errorf("Battle player init failed: %w", err)
	}

	if err := skill.Init(); err != nil {
		return fmt.Errorf("Skill init failed: %w", err)
	}

	if err := enemy.Init(); err != nil {
		return fmt.Errorf("Enemy init failed: %w", err)
	}

	return nil
}

// End ...
func End() {
	field.End()
	battleplayer.End()
	skill.End()
	enemy.End()
}

// Process ...
func Process() error {
	if err := anim.MgrProcess(); err != nil {
		return fmt.Errorf("Failed to handle animation: %w", err)
	}

	switch battleState {
	case stateChipSelect:
		if battleCount == 0 {
			if err := chipsel.Init(battleplayer.Get().ChipFolder); err != nil {
				return fmt.Errorf("Failed to initialize chip select: %w", err)
			}
		}
		if chipsel.Process() {
			// set selected chips
			battleplayer.SetChipSelectResult(chipsel.GetSelected())
			stateChange(stateBeforeMain)
		}
	case stateBeforeMain:
		// TODO implement this
		stateChange(stateMain)
	case stateMain:
		if err := battleplayer.MainProcess(); err != nil {
			if errors.Is(err, battleplayer.ErrChipSelect) {
				stateChange(stateChipSelect)
				return nil
			}
			if errors.Is(err, battleplayer.ErrPlayerDead) {
				// TODO return lose
			}
			return fmt.Errorf("Failed to process player: %w", err)
		}
		if err := enemy.MgrProcess(); err != nil {
			if errors.Is(err, enemy.ErrGameEnd) {
				// TODO return win
			}
			return fmt.Errorf("Failed to process enemy: %w", err)
		}
		fieldUpdates()
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

	switch battleState {
	case stateChipSelect:
		enemy.MgrDraw()
		battleplayer.DrawChar()
		chipsel.Draw()
	case stateMain:
		enemy.MgrDraw()
		battleplayer.DrawChar()
		battleplayer.DrawChipIcon()
	}

	anim.MgrDraw()
}

func fieldUpdates() {
	p := battleplayer.Get()
	objs := []field.ObjectPosition{
		{X: p.PosX, Y: p.PosY, ID: p.ID},
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
