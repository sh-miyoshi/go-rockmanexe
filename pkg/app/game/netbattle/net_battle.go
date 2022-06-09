package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/opening"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/player"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/netconnpb"
)

const (
	stateWaiting int = iota
	stateOpening
	stateChipSelect
	stateWaitSelect
	stateBeforeMain
	stateMain
	stateResult

	stateMax
)

type NetBattle struct {
	conn        *netconn.NetConn
	gameCount   int
	state       int
	stateCount  int
	openingInst opening.Opening
	playerInst  *battleplayer.BattlePlayer
}

var (
	InvalidChips = []int{
		chip.IDBoomerang1,
	}

	inst NetBattle
)

func Init(plyr *player.Player) error {
	logger.Info("Init net battle data ...")
	inst = NetBattle{
		conn:       netconn.GetInst(),
		gameCount:  0,
		state:      stateWaiting,
		stateCount: 0,
	}
	var err error
	inst.openingInst, err = opening.NewWithBoss([]enemy.EnemyParam{
		{CharID: enemy.IDRockman, Pos: common.Point{X: 4, Y: 1}},
	})
	if err != nil {
		return err
	}
	inst.playerInst, err = battleplayer.New(plyr)
	if err != nil {
		return err
	}

	sound.BGMStop()
	return nil
}

func End() {
	inst.conn.Disconnect()
	if inst.openingInst != nil {
		inst.openingInst.End()
	}
}

func Process() error {
	inst.gameCount++

	switch inst.state {
	case stateWaiting:
		status := inst.conn.GetGameStatus()
		if status == pb.Data_CHIPSELECTWAIT {
			if err := sound.BGMPlay(sound.BGMNetBattle); err != nil {
				return fmt.Errorf("failed to play bgm: %v", err)
			}

			stateChange(stateOpening)
			return nil
		}
	case stateOpening:
		if inst.openingInst.Process() {
			stateChange(stateChipSelect)
			return nil
		}
	case stateChipSelect:
		if inst.stateCount == 0 {
			if err := chipsel.Init(inst.playerInst.ChipFolder); err != nil {
				return fmt.Errorf("failed to initialize chip select: %w", err)
			}
		}
		if chipsel.Process() {
			// set selected chips
			// TODO
			stateChange(stateWaitSelect)
			return nil
		}
	case stateWaitSelect:
	}
	// TODO

	inst.stateCount++
	return nil
}

func Draw() {
	switch inst.state {
	case stateWaiting:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		draw.String(140, 110, 0xffffff, "相手の接続を待っています")
	case stateOpening:
		inst.openingInst.Draw()
	case stateChipSelect:
		inst.playerInst.DrawFrame(true, false)
		chipsel.Draw()
	}
}

func stateChange(nextState int) {
	logger.Info("Change battle state from %d to %d", inst.state, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	inst.state = nextState
	inst.stateCount = 0
}
