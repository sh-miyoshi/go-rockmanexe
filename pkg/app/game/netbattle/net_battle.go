package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/opening"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
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
	fieldInst   *field.Field
	drawMgr     *draw.DrawManager
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
	inst.fieldInst, err = field.New()
	if err != nil {
		return err
	}

	inst.drawMgr, err = draw.New(inst.playerInst.Object.ID)
	if err != nil {
		return err
	}

	inst.conn.SendObject(inst.playerInst.Object)
	sound.BGMStop()
	return nil
}

func End() {
	inst.conn.Disconnect()
	if inst.openingInst != nil {
		inst.openingInst.End()
	}
	if inst.fieldInst != nil {
		inst.fieldInst.End()
	}
	if inst.drawMgr != nil {
		inst.drawMgr.End()
	}
}

func Process() error {
	inst.gameCount++
	inst.fieldInst.Update()

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
			inst.playerInst.SetChipSelectResult(chipsel.GetSelected())
			inst.conn.SendObject(inst.playerInst.Object)
			if err := inst.conn.SendSignal(pb.Action_CHIPSEND); err != nil {
				return fmt.Errorf("failed to send Action_CHIPSEND signal: %v", err)
			}
			stateChange(stateWaitSelect)
			return nil
		}
	case stateWaitSelect:
		status := inst.conn.GetGameStatus()
		if status == pb.Data_ACTING {
			stateChange(stateBeforeMain)
			return nil
		}
	case stateBeforeMain:
	case stateMain:
	case stateResult:
	}
	// TODO

	if err := inst.conn.BulkSendData(); err != nil {
		return fmt.Errorf("failed to bulk send data: %w", err)
	}

	inst.stateCount++
	return nil
}

func Draw() {
	inst.fieldInst.Draw()
	inst.drawMgr.DrawObjects()

	switch inst.state {
	case stateWaiting:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		appdraw.String(140, 110, 0xffffff, "相手の接続を待っています")
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
