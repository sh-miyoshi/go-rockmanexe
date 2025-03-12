package netbattle

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	appdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/b4main"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	battlefield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/opening"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
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
	stateCutin

	stateMax
)

type NetBattle struct {
	conn         *netconn.NetConn
	gameCount    int
	state        int
	stateCount   int
	stateInst    battle.State
	playerInst   *battleplayer.BattlePlayer
	fieldInst    *field.Field
	playerInitHP int
	animMgr      *manager.Manager
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
		chip.IDWideShot1,
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
		chip.IDAreaSteal,
		chip.IDBubbleShot,
		chip.IDBubbleSide,
		chip.IDBubbleV,
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
		animMgr:      manager.NewManager(),
	}
	var err error
	inst.stateInst, err = opening.NewWithBoss([]enemy.EnemyParam{
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
		return errors.Wrap(err, "effect init failed")
	}

	sound.BGMStop()
	return nil
}

func End() {
	inst.conn.Disconnect()
	if inst.stateInst != nil {
		inst.stateInst.End()
		inst.stateInst = nil
	}
	if inst.fieldInst != nil {
		inst.fieldInst.End()
	}
	draw.End()
	inst.animMgr.AnimCleanup()
}

func Update() error {
	battlecommon.SystemProcess()
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
				return errors.Wrap(err, "failed to send initial player param")
			}

			if err := sound.BGMPlay(sound.BGMNetBattle); err != nil {
				return errors.Wrap(err, "failed to play bgm")
			}

			stateChange(stateOpening)
			return nil
		}
	case stateOpening:
		if inst.stateInst.Update() {
			inst.stateInst.End()
			inst.stateInst = nil
			stateChange(stateChipSelect)
			return nil
		}
	case stateChipSelect:
		if inst.stateCount == 0 {
			if err := chipsel.Init(inst.playerInst.GetChipFolder(), inst.playerInst.GetChipSelectMax()); err != nil {
				return errors.Wrap(err, "failed to initialize chip select")
			}
		}
		if chipsel.Update() {
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
			inst.stateInst, err = b4main.New(inst.playerInst.GetSelectedChips())
			if err != nil {
				return errors.Wrap(err, "failed to initialize before main")
			}
			inst.playerInst.UpdatePA()
		}

		if inst.stateInst.Update() {
			inst.stateInst.End()
			inst.stateInst = nil
			stateChange(stateMain)
			return nil
		}
	case stateMain:
		isRunAnim = true
		done, err := inst.playerInst.Update()
		if err != nil {
			return errors.Wrap(err, "player process failed")
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
		case pb.Response_CUTIN:
			stateChange(stateCutin)
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
			inst.stateInst, err = titlemsg.New(fname, 60)
			if err != nil {
				return errors.Wrap(err, "failed to initialize result")
			}
		}

		if inst.stateInst.Update() {
			inst.stateInst.End()
			inst.stateInst = nil
			if inst.playerInst.IsDead() {
				return battle.ErrLose
			}
			return battle.ErrWin
		}
	case stateCutin:
		isRunAnim = true
		// 待っている状態
		if inst.stateCount == 0 {
			battlefield.SetBlackoutCount(9999) // 正確なデータが得られるまで一旦セットしておく
		}

		// 暗転チップの情報を処理
		cutin, isSet := inst.conn.PopCutinInfo()
		if isSet {
			clientID := config.Get().Net.ClientID
			skill.SetChipNameDraw(cutin.SkillName, clientID == cutin.OwnerClientID)
		}

		status := inst.conn.GetGameStatus()
		if status == pb.Response_ACTING {
			battlefield.SetBlackoutCount(0)
			stateChange(stateMain)
			return nil
		}
	}

	if isRunAnim {
		if err := inst.animMgr.AnimMgrProcess(); err != nil {
			return errors.Wrap(err, "failed to handle animation")
		}
	}

	inst.stateCount++
	return nil
}

func Draw() {
	inst.fieldInst.Draw()
	draw.Draw()

	inst.animMgr.AnimMgrDraw()
	inst.playerInst.LocalDraw()

	battlefield.DrawBlackout()

	switch inst.state {
	case stateWaiting:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, config.ScreenSize.X, config.ScreenSize.Y, 0, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		appdraw.String(140, 110, 0xffffff, "相手の接続を待っています")
	case stateOpening:
		if inst.stateInst != nil {
			inst.stateInst.Draw()
		}
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
		if inst.stateInst != nil {
			inst.stateInst.Draw()
		}
	case stateMain:
		inst.playerInst.DrawFrame(false, true)
	case stateResult:
		inst.playerInst.DrawFrame(false, true)
		if inst.stateInst != nil {
			inst.stateInst.Draw()
		}
	}

	battlecommon.SystemDraw()
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
		inst.animMgr.EffectAnimNew(effect.Get(e.Type, e.Pos, e.RandRange))
	}
}

func handleSound() {
	inst := net.GetInst()
	for _, s := range inst.GetGameInfo().Sounds {
		sound.On(resources.SEType(s.Type))
	}
	inst.CleanupSounds()
}
