package battle

import (
	"fmt"

	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/b4main"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/chipsel"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/enemy"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/opening"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/win"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/win/reward"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateOpening int = iota
	stateChipSelect
	stateBeforeMain
	stateMain
	stateResultWin
	stateResultLose

	stateMax
)

type State interface {
	End()
	Update() bool
	Draw()
}

var (
	isStateInit    bool
	battleState    int
	playerInst     *battleplayer.BattlePlayer
	enemyList      []enemy.EnemyParam
	gameCount      int
	stateInst      State
	basePlayerInst *player.Player
	animMgr        *manager.Manager
	chipSelect     *chipsel.ChipSelect

	ErrWin  = errors.New("player win")
	ErrLose = errors.New("playser lose")
)

func Init(plyr *player.Player, enemies []enemy.EnemyParam) error {
	logger.Info("Init battle data ...")

	gameCount = 0
	isStateInit = false
	battleState = stateOpening
	basePlayerInst = plyr
	stateInst = nil
	animMgr = manager.NewManager()

	var err error
	playerInst, err = battleplayer.New(plyr, animMgr)
	if err != nil {
		return errors.Wrap(err, "battle player init failed")
	}
	animMgr.ObjAnimNew(playerInst)
	chipSelect = chipsel.NewChipSelect(playerInst)

	enemyList = []enemy.EnemyParam{}
	for _, e := range enemies {
		if e.CharID == enemy.IDSupportNPC {
			// Supporterは必ずPlayerを作成後に作成する
			supporter, err := battleplayer.NewSupporter(battleplayer.SupporterParam{
				HP:      uint(e.HP),
				InitPos: e.Pos,
			}, animMgr)
			if err != nil {
				return errors.Wrap(err, "battle supporter init failed")
			}
			animMgr.ObjAnimNew(supporter)
			logger.Info("add supporter %+v", supporter)
		} else {
			enemyList = append(enemyList, e)
		}
	}
	logger.Info("enemy list: %+v", enemyList)

	if err := field.Init(animMgr); err != nil {
		return errors.Wrap(err, "battle field init failed")
	}

	if err := skill.Init(); err != nil {
		return errors.Wrap(err, "skill init failed")
	}

	if err := effect.Init(); err != nil {
		return errors.Wrap(err, "effect init failed")
	}

	bgm := sound.BGMBattle
	for _, e := range enemies {
		if enemy.IsBoss(e.CharID) {
			bgm = sound.BGMBossBattle
			break
		}
	}

	if err := sound.BGMPlay(bgm); err != nil {
		return errors.Wrap(err, "failed to play bgm")
	}

	// 敵討伐時の報酬をセット
	reward.SetEnemyChipList([]reward.EnemyChipInfo{
		{CharID: enemy.IDMetall, ChipID: chip.IDShockWave, Code: "l", RequiredLevel: 7},
		{CharID: enemy.IDMetall, ChipID: chip.IDShockWave, Code: "*", RequiredLevel: 9},
		{CharID: enemy.IDBilly, ChipID: chip.IDThunderBall1, Code: "l", RequiredLevel: 7},
		{CharID: enemy.IDLark, ChipID: chip.IDWideShot1, Code: "c", RequiredLevel: 7},
		{CharID: enemy.IDBoomer, ChipID: chip.IDBoomerang1, Code: "m", RequiredLevel: 7},
		{CharID: enemy.IDBoomer, ChipID: chip.IDBoomerang1, Code: "*", RequiredLevel: 9},
		{CharID: enemy.IDAquaman, ChipID: chip.IDAquaman, Code: "a", RequiredLevel: 9},
		{CharID: enemy.IDVolgear, ChipID: chip.IDFlameLine1, Code: "f", RequiredLevel: 7},
		{CharID: enemy.IDGaroo, ChipID: chip.IDHeatShot, Code: "c", RequiredLevel: 7},
		{CharID: enemy.IDShrimpy, ChipID: chip.IDBubbleShot, Code: "b", RequiredLevel: 7},
		{CharID: enemy.IDShrimpy, ChipID: chip.IDBubbleShot, Code: "c", RequiredLevel: 7},
		{CharID: enemy.IDShrimpy, ChipID: chip.IDBubbleSide, Code: "f", RequiredLevel: 9},
		{CharID: enemy.IDShrimpy, ChipID: chip.IDBubbleV, Code: "f", RequiredLevel: 9},
		{CharID: enemy.IDForte, ChipID: chip.IDForteAnother, Code: "x", RequiredLevel: 1, IsOnlyOne: true},
		// TODO: コールドマン、サーキラーのチップ
	})

	// カスタムゲージのスピードをデフォルトにしておく
	battlecommon.CustomGaugeSpeed = battlecommon.DefaultCustomGaugeSpeed

	logger.Info("Successfully initialized battle data")
	return nil
}

func End() {
	field.ResetSet4x4Area()
	animMgr.Cleanup()
	field.End()
	playerInst.End()
	skill.End()
	enemy.End()
	effect.End()
	if chipSelect != nil {
		chipSelect.End()
	}
	logger.Info("End battle data")
}

func Update() error {
	battlecommon.SystemProcess()

	isActive := false

	switch battleState {
	case stateOpening:
		if !isStateInit {
			isStateInit = true
			var err error
			if enemy.IsBoss(enemyList[0].CharID) {
				stateInst, err = opening.NewWithBoss(enemyList)
			} else {
				stateInst, err = opening.NewWithNormal(enemyList)
			}
			if err != nil {
				return errors.Wrap(err, "opening init failed")
			}
		}

		if stateInst.Update() {
			if err := enemy.Init(playerInst.ID, enemyList, animMgr); err != nil {
				return errors.Wrap(err, "enemy init failed")
			}
			stateChange(stateChipSelect)
			return nil
		}
	case stateChipSelect:
		if !isStateInit {
			isStateInit = true
			if err := chipSelect.Init(playerInst.ChipFolder, playerInst.ChipSelectMax); err != nil {
				return errors.Wrap(err, "failed to initialize chip select")
			}
			playerInst.SetFrameInfo(true, false)
		}
		if chipSelect.Update() {
			// set selected chips
			playerInst.SetChipSelectResult(chipSelect.GetSelected())
			stateChange(stateBeforeMain)
			return nil
		}
	case stateBeforeMain:
		if !isStateInit {
			isStateInit = true
			var err error
			stateInst, err = b4main.New(playerInst.SelectedChips)
			if err != nil {
				return errors.Wrap(err, "failed to initialize before main")
			}
			playerInst.UpdateChipInfo()
			playerInst.SetFrameInfo(false, true)
			playerInst.UpdateStatus()
		}

		if stateInst.Update() {
			stateChange(stateMain)
			return nil
		}
	case stateMain:
		gameCount++

		isActive = !field.IsBlackout()

		switch playerInst.NextAction {
		case battleplayer.NextActChipSelect:
			stateChange(stateChipSelect)
			playerInst.NextAction = battleplayer.NextActNone
			return nil
		case battleplayer.NextActLose:
			cleanupBattleAnims()
			stateChange(stateResultLose)
			return nil
		}
		if err := enemy.MgrProcess(); err != nil {
			if errors.Is(err, enemy.ErrGameEnd) {
				cleanupBattleAnims()
				playerInst.EnableAct = false
				stateChange(stateResultWin)
				return nil
			}
			return errors.Wrap(err, "failed to process enemy")
		}

		field.Update()
	case stateResultWin:
		if !isStateInit {
			isStateInit = true
			enemies := []reward.EnemyParam{}
			for _, e := range enemyList {
				enemies = append(enemies, reward.EnemyParam{
					CharID: e.CharID,
					IsBoss: enemy.IsBoss(e.CharID),
				})
			}

			var err error
			stateInst, err = win.New(reward.WinArg{
				GameTime:        gameCount,
				DeletedEnemies:  enemies,
				PlayerMoveNum:   playerInst.MoveNum,
				PlayerDamageNum: playerInst.DamageNum,
			}, basePlayerInst)
			if err != nil {
				return errors.Wrap(err, "failed to initialize result win")
			}
		}

		if stateInst.Update() {
			stateInst.End()
			stateInst = nil
			return ErrWin
		}
	case stateResultLose:
		if !isStateInit {
			isStateInit = true
			fname := config.ImagePath + "battle/msg_lose.png"
			var err error
			stateInst, err = titlemsg.New(fname, 0)
			if err != nil {
				return errors.Wrap(err, "failed to initialize lose")
			}
			playerInst.SetFrameInfo(false, false)
		}

		if stateInst.Update() {
			sound.SEClear()
			return ErrLose
		}
	}

	if err := animMgr.Update(isActive); err != nil {
		return errors.Wrap(err, "failed to handle animation")
	}

	return nil
}

func Draw() {
	field.Draw()
	animMgr.Draw()

	drawEnemyNames()
	field.DrawBlackout()

	if stateInst != nil {
		stateInst.Draw()
	}

	switch battleState {
	case stateChipSelect:
		chipSelect.Draw()
	}

	battlecommon.SystemDraw()
}

func stateChange(nextState int) {
	logger.Info("Change battle state from %d to %d", battleState, nextState)
	if nextState < 0 || nextState >= stateMax {
		system.SetError(fmt.Sprintf("Invalid next battle state: %d", nextState))
		return
	}
	battleState = nextState
	isStateInit = false
	if stateInst != nil {
		stateInst.End()
		stateInst = nil
	}
}

func drawEnemyNames() {
	for i, e := range enemyList {
		name := enemy.GetName(e.CharID)
		ofs := dxlib.GetDrawStringWidth(name, len(name))
		draw.String(config.ScreenSize.X-ofs-5, i*20+10, 0xffffff, "%s", name)
	}
}

func cleanupBattleAnims() {
	for _, a := range animMgr.AnimGetSkills() {
		animMgr.AnimDelete(a.ObjID)
	}
	animMgr.DamageManager().RemoveAll()
}
