package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
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
	logger.Info("Init net battle data ...")

	gameCount = 0
	battleCount = 0
	battleState = stateWaiting
	b4mainInst = nil
	loseInst = nil
	basePlayerInst = plyr

	f, err := netconn.GetFieldInfo()
	if err != nil {
		return fmt.Errorf("get field info failed: %w", err)
	}
	field.Init(f)

	// TODO

	return nil
}

func End() {
	netconn.Disconnect()
}

func Process() error {
	switch battleState {
	case stateWaiting:
		status, err := netconn.GetStatus()
		if err != nil {
			return fmt.Errorf("get connect status error: %w", err)
		}
		if status == pb.Data_CHIPSELECTWAIT {
			battleState = stateOpening
		}
	case stateOpening:
		// TODO
	}

	battleCount++
	return nil
}

func Draw() {
	field.Draw()
}

func stateChange(nextState int) {
	logger.Info("Change battle state from %d to %d", battleState, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	battleState = nextState
	battleCount = 0
}
