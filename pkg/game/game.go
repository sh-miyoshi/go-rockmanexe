package game

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

const (
	stateTitle int = iota
	stateBattle
	stateMax
)

var (
	state           = stateTitle
	count      uint = 0
	playerInfo *player.Player
)

// Process ...
func Process() error {
	switch state {
	case stateTitle:
		// TODO
		// show opening page
		// select "はじめから" or "つづきから"
		playerInfo = player.New()
		stateChange(stateBattle)
		return nil
	case stateBattle:
		if count == 0 {
			if err := battle.Init(playerInfo); err != nil {
				return fmt.Errorf("Game process in state battle failed: %w", err)
			}
		}
		if err := battle.Process(); err != nil {
			return fmt.Errorf("Battle process failed: % w", err)
		}
	}
	count++
	return nil
}

// Draw ...
func Draw() {
	if count == 0 {
		// skip if initialize phase
		return
	}

	switch state {
	case stateBattle:
		battle.Draw()
	}
}

func stateChange(nextState int) {
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next game state: %d", nextState))
	}
	state = nextState
	count = 0
}
