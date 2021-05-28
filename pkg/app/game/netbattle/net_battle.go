package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/routerpb"
)

const (
	stateWaiting int = iota
	stateOpening
	stateChipSelect
	stateWaitSelect
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
	playerInst     *battleplayer.BattlePlayer
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

	playerInst, err = battleplayer.New(plyr)
	if err != nil {
		return fmt.Errorf("net battle player init failed: %w", err)
	}

	// TODO

	return nil
}

func End() {
	netconn.Disconnect()
	playerInst.End()
}

func Process() error {
	switch battleState {
	case stateWaiting:
		status, err := netconn.GetStatus()
		if err != nil {
			return fmt.Errorf("get connect status error: %w", err)
		}
		if status == pb.Data_CHIPSELECTWAIT {
			stateChange(stateOpening)
			return nil
		}
	case stateOpening:
		// TODO animation処理
		stateChange(stateChipSelect)
		return nil
	case stateChipSelect:
		if battleCount == 0 {
			if err := chipsel.Init(playerInst.ChipFolder); err != nil {
				return fmt.Errorf("failed to initialize chip select: %w", err)
			}
		}
		if chipsel.Process() {
			// set selected chips
			playerInst.SetChipSelectResult(chipsel.GetSelected())
			// TODO send to server
			stateChange(stateWaitSelect)
			return nil
		}
	case stateWaitSelect:
		// TODO
	}

	battleCount++
	return nil
}

func Draw() {
	field.Draw()

	switch battleState {
	case stateWaiting:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenX, common.ScreenY, 0, dxlib.TRUE)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		draw.String(140, 110, 0xffffff, "相手の接続を待っています")
	case stateOpening:
		// TODO animation
	case stateChipSelect:
		playerInst.DrawFrame(true, false)
		chipsel.Draw()
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
