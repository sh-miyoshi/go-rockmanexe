package netbattle

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	netdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/opening"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
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
	stateResult

	stateMax
)

var (
	battleCount int
	battleState int
	gameCount   int
	b4mainInst  *titlemsg.TitleMsg
	resultInst  *titlemsg.TitleMsg
	playerInst  *battleplayer.BattlePlayer

	InvalidChips = []int{
		chip.IDBoomerang1,
	}
)

func Init(plyr *player.Player) error {
	logger.Info("Init net battle data ...")

	gameCount = 0
	battleCount = 0
	battleState = stateWaiting
	b4mainInst = nil
	resultInst = nil

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

	netconn.SendObject(playerInst.Object)

	sound.BGMStop()
	skill.Init(playerInst.Object.ID)

	return nil
}

func End() {
	netconn.Disconnect()
	playerInst.End()
	netdraw.End()
}

func Process() error {
	effect.Process()
	soundProc()
	status := netconn.GetStatus()

	switch battleState {
	case stateWaiting:
		if status == pb.Data_CHIPSELECTWAIT {
			if err := sound.BGMPlay(sound.BGMNetBattle); err != nil {
				return fmt.Errorf("failed to play bgm: %v", err)
			}

			// Delay udpate due to Cloud Run Rate Exceeded
			if err := netconn.BulkSendFieldInfo(); err != nil {
				return fmt.Errorf("failed to add init player object: %w", err)
			}

			stateChange(stateOpening)
			return nil
		}
	case stateOpening:
		if battleCount == 0 {
			opening.Init()
		}
		if opening.Process() {
			stateChange(stateChipSelect)
			return nil
		}
	case stateChipSelect:
		if battleCount == 0 {
			if err := chipsel.Init(playerInst.ChipFolder); err != nil {
				return fmt.Errorf("failed to initialize chip select: %w", err)
			}
		}
		if chipsel.Process() {
			// set selected chips
			playerInst.SetChipSelectResult(chipsel.GetSelected())
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
			playerInst.InitBattleFrame()
			fname := common.ImagePath + "battle/msg_start.png"
			var err error
			b4mainInst, err = titlemsg.New(fname, 0)
			if err != nil {
				return fmt.Errorf("failed to initialize before main: %w", err)
			}
		}

		netconn.UpdateObjectsCount()
		if err := netconn.BulkSendFieldInfo(); err != nil {
			return fmt.Errorf("send field info failed: %w", err)
		}

		if b4mainInst.Process() {
			b4mainInst.End()
			stateChange(stateMain)
			return nil
		}
	case stateMain:
		gameCount++
		netconn.UpdateObjectsCount()
		if err := netconn.BulkSendFieldInfo(); err != nil {
			return fmt.Errorf("send field info failed: %w", err)
		}

		if status == pb.Data_CHIPSELECTWAIT {
			stateChange(stateChipSelect)
			return nil
		} else if status == pb.Data_GAMEEND {
			stateChange(stateResult)
			return nil
		}

		done, err := playerInst.Process()
		if err != nil {
			return fmt.Errorf("player process failed: %w", err)
		}
		if done {
			stateChange(stateResult)
			return nil
		}

		if err := skill.Process(); err != nil {
			return fmt.Errorf("skill process failed: %w", err)
		}
	case stateResult:
		if battleCount == 0 {
			netconn.Disconnect()

			fname := common.ImagePath + "battle/msg_win.png"
			if playerInst.Object.HP <= 0 {
				fname = common.ImagePath + "battle/msg_lose.png"
			}

			var err error
			resultInst, err = titlemsg.New(fname, 60)
			if err != nil {
				return fmt.Errorf("failed to initialize result: %w", err)
			}
		}

		if resultInst.Process() {
			resultInst.End()
			if playerInst.Object.HP <= 0 {
				return battle.ErrLose
			}
			return battle.ErrWin
		}
	}

	battleCount++
	return nil
}

func Draw() {
	field.Draw(playerInst.Object.ID)
	playerInst.DrawOptions()

	switch battleState {
	case stateWaiting:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenX, common.ScreenY, 0, dxlib.TRUE)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
		draw.String(140, 110, 0xffffff, "相手の接続を待っています")
	case stateOpening:
		opening.Draw()
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
	case stateResult:
		if resultInst != nil {
			resultInst.Draw()
		}
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

func soundProc() {
	finfo := netconn.GetFieldInfo()
	for _, se := range finfo.Sounds {
		sound.On(sound.SEType(se))
	}
	netconn.RemoveSounds()
}
