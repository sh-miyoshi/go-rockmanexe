package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
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
	conn       *netconn.NetConn
	gameCount  int
	state      int
	stateCount int
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

	sound.BGMStop()
	return nil
}

func End() {
	inst.conn.Disconnect()
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
		if inst.stateCount == 0 {

		}
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
