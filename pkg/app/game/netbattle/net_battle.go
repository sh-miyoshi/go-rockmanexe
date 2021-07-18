package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	netdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/skill"
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
	battleCount    int
	battleState    int
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

	if err := field.Init(); err != nil {
		return fmt.Errorf("net battle field init failed: %w", err)
	}

	var err error
	playerInst, err = battleplayer.New(plyr)
	if err != nil {
		return fmt.Errorf("net battle player init failed: %w", err)
	}

	if err := netdraw.Init(); err != nil {
		return fmt.Errorf("failed to init battle draw info: %w", err)
	}

	if err := netconn.SendObject(playerInst.Object); err != nil {
		return fmt.Errorf("failed to add init player object: %w", err)
	}

	// TODO

	return nil
}

func End() {
	netconn.Disconnect()
	playerInst.End()
	netdraw.End()
}

func Process() error {
	status, err := netconn.GetStatus()
	if err != nil {
		return fmt.Errorf("get connect status error: %w", err)
	}
	netconn.UpdateObjectsCount()
	effect.Process()

	switch battleState {
	case stateWaiting:
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
			// TODO error handling
			netconn.SendObject(playerInst.Object)
			netconn.SendSignal(pb.Action_CHIPSEND)
			stateChange(stateWaitSelect)
			return nil
		}
	case stateWaitSelect:
		if status == pb.Data_ACTING {
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

		if status == pb.Data_CHIPSELECTWAIT {
			stateChange(stateChipSelect)
			return nil
		}

		// TODO gameend, error handling
		playerInst.Process()

		if err := skill.Process(); err != nil {
			return fmt.Errorf("skill process failed: %w", err)
		}
	}

	battleCount++
	return nil
}

func Draw() {
	field.Draw(playerInst.Object.ID)
	playerInst.DrawChargeShot()

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
	case stateWaitSelect:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenX, common.ScreenY, 0, dxlib.TRUE)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		draw.String(140, 110, 0xffffff, "相手の選択を待っています")
	case stateBeforeMain:
		playerInst.DrawFrame(false, true)
		if b4mainInst != nil {
			b4mainInst.Draw()
		}
	case stateMain:
		playerInst.DrawFrame(false, true)
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
