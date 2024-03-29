package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	appdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/b4main"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/opening"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	pb "github.com/sh-miyoshi/go-rockmanexe/pkg/net/netconnpb"
	netobj "github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
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
	conn         *netconn.NetConn
	gameCount    int
	state        int
	stateCount   int
	openingInst  opening.Opening
	playerInst   *battleplayer.BattlePlayer
	fieldInst    *field.Field
	b4mainInst   *b4main.BeforeMain
	resultInst   *titlemsg.TitleMsg
	playerInitHP int
}

var (
	ValidChips = []int{
		chip.IDCannon,
		chip.IDHighCannon,
		chip.IDMegaCannon,
		chip.IDMiniBomb,
		chip.IDRecover10,
		chip.IDRecover30,
		chip.IDRecover50,
		chip.IDRecover80,
		chip.IDRecover120,
		chip.IDRecover150,
		chip.IDRecover200,
		chip.IDRecover300,
		chip.IDShockWave,
		chip.IDSpreadGun,
		chip.IDSword,
		chip.IDWideSword,
		chip.IDLongSword,
		chip.IDVulcan1,
		chip.IDWideShot,
		chip.IDHeatShot,
		chip.IDHeatV,
		chip.IDHeatSide,
		chip.IDFlameLine1,
		chip.IDFlameLine2,
		chip.IDFlameLine3,
		chip.IDTornado,
		chip.IDBoomerang1,
		chip.IDBambooLance,
		chip.IDCrackout,
		chip.IDDoubleCrack,
		chip.IDTripleCrack,
	}

	inst NetBattle
)

func Init(plyr *player.Player) error {
	logger.Info("Init net battle data ...")
	inst = NetBattle{
		conn:         net.GetInst(),
		gameCount:    0,
		state:        stateWaiting,
		stateCount:   0,
		playerInitHP: int(plyr.HP),
	}
	var err error
	inst.openingInst, err = opening.NewWithBoss([]enemy.EnemyParam{
		{CharID: enemy.IDRockman, Pos: point.Point{X: 4, Y: 1}},
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

	if err = draw.Init(inst.playerInst.GetObjectID()); err != nil {
		return err
	}

	if err := effect.Init(); err != nil {
		return fmt.Errorf("effect init failed: %w", err)
	}

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
	draw.End()
	localanim.AnimCleanup()
}

func Process() error {
	inst.gameCount++
	isRunAnim := false

	handleSound()

	switch inst.state {
	case stateWaiting:
		status := inst.conn.GetGameStatus()
		if status == pb.Response_CHIPSELECTWAIT {
			obj := netobj.InitParam{
				ID: inst.playerInst.GetObjectID(),
				HP: inst.playerInitHP,
				X:  1,
				Y:  1,
			}
			if err := inst.conn.SendSignal(pb.Request_INITPARAMS, obj.Marshal()); err != nil {
				return fmt.Errorf("failed to send initial player param: %w", err)
			}

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
			if err := chipsel.Init(inst.playerInst.GetChipFolder()); err != nil {
				return fmt.Errorf("failed to initialize chip select: %w", err)
			}
		}
		if chipsel.Process() {
			// set selected chips
			inst.playerInst.SetChipSelectResult(chipsel.GetSelected())
			// TODO: 選択したチップ一覧を送る
			inst.conn.SendSignal(pb.Request_CHIPSELECT, nil)
			stateChange(stateWaitSelect)
			return nil
		}
	case stateWaitSelect:
		status := inst.conn.GetGameStatus()
		if status == pb.Response_ACTING {
			stateChange(stateBeforeMain)
			return nil
		}
	case stateBeforeMain:
		if inst.stateCount == 0 {
			var err error
			inst.b4mainInst, err = b4main.New(inst.playerInst.GetSelectedChips())
			if err != nil {
				return fmt.Errorf("failed to initialize before main: %w", err)
			}
			inst.playerInst.UpdatePA()
		}

		if inst.b4mainInst.Process() {
			inst.b4mainInst.End()
			stateChange(stateMain)
			return nil
		}
	case stateMain:
		isRunAnim = true
		done, err := inst.playerInst.Process()
		if err != nil {
			return fmt.Errorf("player process failed: %w", err)
		}
		if done {
			stateChange(stateResult)
			return nil
		}

		handleEffect()

		status := inst.conn.GetGameStatus()
		switch status {
		case pb.Response_CHIPSELECTWAIT:
			stateChange(stateChipSelect)
			return nil
		case pb.Response_GAMEEND:
			stateChange(stateResult)
			return nil
		}
	case stateResult:
		isRunAnim = true
		if inst.stateCount == 0 {
			net.GetInst().Disconnect()

			fname := config.ImagePath + "battle/msg_win.png"
			if inst.playerInst.IsDead() {
				fname = config.ImagePath + "battle/msg_lose.png"
			}

			var err error
			inst.resultInst, err = titlemsg.New(fname, 60)
			if err != nil {
				return fmt.Errorf("failed to initialize result: %w", err)
			}
		}

		if inst.resultInst.Process() {
			inst.resultInst.End()
			if inst.playerInst.IsDead() {
				return battle.ErrLose
			}
			return battle.ErrWin
		}
	}

	if isRunAnim {
		// TODO(blackout中はエフェクトもとめておく？)
		if err := localanim.AnimMgrProcess(); err != nil {
			return fmt.Errorf("failed to handle animation: %w", err)
		}
	}

	inst.stateCount++
	return nil
}

func Draw() {
	inst.fieldInst.Draw()
	draw.Draw()

	localanim.AnimMgrDraw()
	inst.playerInst.LocalDraw()

	switch inst.state {
	case stateWaiting:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, config.ScreenSize.X, config.ScreenSize.Y, 0, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		appdraw.String(140, 110, 0xffffff, "相手の接続を待っています")
	case stateOpening:
		inst.openingInst.Draw()
	case stateChipSelect:
		inst.playerInst.DrawFrame(true, false)
		chipsel.Draw()
	case stateWaitSelect:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, config.ScreenSize.X, config.ScreenSize.Y, 0, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		appdraw.String(140, 110, 0xffffff, "相手の選択を待っています")
	case stateBeforeMain:
		inst.playerInst.DrawFrame(false, true)
		if inst.b4mainInst != nil {
			inst.b4mainInst.Draw()
		}
	case stateMain:
		inst.playerInst.DrawFrame(false, true)
	case stateResult:
		inst.playerInst.DrawFrame(false, true)
		if inst.resultInst != nil {
			inst.resultInst.Draw()
		}
	}
}

func stateChange(nextState int) {
	logger.Info("Change battle state from %d to %d", inst.state, nextState)
	if nextState < 0 || nextState >= stateMax {
		system.SetError(fmt.Sprintf("Invalid next battle state: %d", nextState))
	}
	inst.state = nextState
	inst.stateCount = 0
}

func handleEffect() {
	g := net.GetInst().GetGameInfo()
	for _, e := range g.Effects {
		localanim.AnimNew(effect.Get(e.Type, e.Pos, e.RandRange))
	}
}

func handleSound() {
	inst := net.GetInst()
	for _, s := range inst.GetGameInfo().Sounds {
		sound.On(resources.SEType(s.Type))
	}
	inst.CleanupSounds()
}
