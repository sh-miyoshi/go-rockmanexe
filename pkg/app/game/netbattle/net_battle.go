package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateWaiting int = iota
	stateOpening
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
	// playerInst     *battleplayer.BattlePlayer
	gameCount      int
	b4mainInst     *titlemsg.TitleMsg
	loseInst       *titlemsg.TitleMsg
	basePlayerInst *player.Player
)

func Init(plyr *player.Player) error {
	logger.Info("Init battle data ...")

	gameCount = 0
	battleCount = 0
	battleState = stateOpening
	b4mainInst = nil
	loseInst = nil
	basePlayerInst = plyr

	// TODO

	return nil
}

func End() {

}

func Process() error {
	switch battleState {
	}

	battleCount++
	return nil
}

func Draw() {
	switch battleState {
	}
}

func stateChange(nextState int) {
	logger.Info("Change battle state from %d to %d", battleState, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	battleState = nextState
	battleCount = 0
}
