package battle

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
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

var (
	battleCount    int
	battleState    int
	playerInst     *battleplayer.BattlePlayer
	enemyList      []enemy.EnemyParam
	gameCount      int
	b4mainInst     *b4main.BeforeMain
	loseInst       *titlemsg.TitleMsg
	openingInst    opening.Opening
	basePlayerInst *player.Player

	ErrWin  = errors.New("player win")
	ErrLose = errors.New("playser lose")
)

func Init(plyr *player.Player, enemies []enemy.EnemyParam) error {
	logger.Info("Init battle data ...")

	gameCount = 0
	battleCount = 0
	battleState = stateOpening
	b4mainInst = nil
	loseInst = nil
	basePlayerInst = plyr

	var err error
	playerInst, err = battleplayer.New(plyr)
	if err != nil {
		return fmt.Errorf("battle player init failed: %w", err)
	}
	localanim.ObjAnimNew(playerInst)

	enemyList = []enemy.EnemyParam{}
	for _, e := range enemies {
		if e.CharID == enemy.IDSupportNPC {
			// Supporterは必ずPlayerを作成後に作成する
			supporter, err := battleplayer.NewSupporter(battleplayer.SupporterParam{
				HP:      uint(e.HP),
				InitPos: e.Pos,
			})
			if err != nil {
				return fmt.Errorf("battle supporter init failed: %w", err)
			}
			localanim.ObjAnimNew(supporter)
			logger.Info("add supporter %+v", supporter)
		} else {
			enemyList = append(enemyList, e)
		}
	}
	logger.Info("enemy list: %+v", enemyList)

	if err := field.Init(); err != nil {
		return fmt.Errorf("battle field init failed: %w", err)
	}

	if err := skill.Init(); err != nil {
		return fmt.Errorf("skill init failed: %w", err)
	}

	if err := effect.Init(); err != nil {
		return fmt.Errorf("effect init failed: %w", err)
	}

	bgm := sound.BGMBattle
	for _, e := range enemies {
		if enemy.IsBoss(e.CharID) {
			bgm = sound.BGMBossBattle
			break
		}
	}

	if err := sound.BGMPlay(bgm); err != nil {
		return fmt.Errorf("failed to play bgm: %v", err)
	}

	// カスタムゲージのスピードをデフォルトにしておく
	battlecommon.CustomGaugeSpeed = battlecommon.DefaultCustomGaugeSpeed

	logger.Info("Successfully initialized battle data")
	return nil
}

func End() {
	field.ResetSet4x4Area()
	localanim.AnimCleanup()
	localanim.ObjAnimCleanup()
	field.End()
	playerInst.End()
	skill.End()
	enemy.End()
	effect.End()
	win.End()
	logger.Info("End battle data")
}

func Process() error {
	battlecommon.SystemProcess()
	isRunAnim := false

	switch battleState {
	case stateOpening:
		if battleCount == 0 {
			var err error
			if enemy.IsBoss(enemyList[0].CharID) {
				openingInst, err = opening.NewWithBoss(enemyList)
			} else {
				openingInst, err = opening.NewWithNormal(enemyList)
			}
			if err != nil {
				return fmt.Errorf("opening init failed: %w", err)
			}
		}

		if openingInst.Process() {
			openingInst.End()
			if err := enemy.Init(playerInst.ID, enemyList); err != nil {
				return fmt.Errorf("enemy init failed: %w", err)
			}
			stateChange(stateChipSelect)
			return nil
		}
	case stateChipSelect:
		if battleCount == 0 {
			if err := chipsel.Init(playerInst.ChipFolder, playerInst.ChipSelectMax); err != nil {
				return fmt.Errorf("failed to initialize chip select: %w", err)
			}
			playerInst.SetFrameInfo(true, false)
		}
		if chipsel.Process() {
			// set selected chips
			playerInst.SetChipSelectResult(chipsel.GetSelected())
			stateChange(stateBeforeMain)
			return nil
		}
	case stateBeforeMain:
		if battleCount == 0 {
			var err error
			b4mainInst, err = b4main.New(playerInst.SelectedChips)
			if err != nil {
				return fmt.Errorf("failed to initialize before main: %w", err)
			}
			playerInst.UpdateChipInfo()
			playerInst.SetFrameInfo(false, true)
		}

		if b4mainInst.Process() {
			b4mainInst.End()
			stateChange(stateMain)
			return nil
		}
	case stateMain:
		isRunAnim = true
		gameCount++

		if err := localanim.ObjAnimMgrProcess(true, field.IsBlackout()); err != nil {
			return fmt.Errorf("failed to handle object animation: %w", err)
		}

		if !field.IsBlackout() {
			switch playerInst.NextAction {
			case battleplayer.NextActChipSelect:
				stateChange(stateChipSelect)
				playerInst.NextAction = battleplayer.NextActNone
				return nil
			case battleplayer.NextActLose:
				stateChange(stateResultLose)
				return nil
			}
			if err := enemy.MgrProcess(); err != nil {
				if errors.Is(err, enemy.ErrGameEnd) {
					playerInst.EnableAct = false
					stateChange(stateResultWin)
					return nil
				}
				return fmt.Errorf("failed to process enemy: %w", err)
			}
		}

		field.Update()
	case stateResultWin:
		isRunAnim = true
		if battleCount == 0 {
			if err := win.Init(win.WinArg{
				GameTime:        gameCount,
				DeletedEnemies:  enemyList,
				PlayerMoveNum:   playerInst.MoveNum,
				PlayerDamageNum: playerInst.DamageNum,
			}, basePlayerInst); err != nil {
				return fmt.Errorf("failed to initialize result win: %w", err)
			}
		}

		if err := localanim.ObjAnimMgrProcess(false, field.IsBlackout()); err != nil {
			return fmt.Errorf("failed to handle object animation: %w", err)
		}

		if win.Process() {
			return ErrWin
		}
	case stateResultLose:
		isRunAnim = true
		if battleCount == 0 {
			fname := config.ImagePath + "battle/msg_lose.png"
			var err error
			loseInst, err = titlemsg.New(fname, 0)
			if err != nil {
				return fmt.Errorf("failed to initialize lose: %w", err)
			}
			playerInst.SetFrameInfo(false, false)
		}

		if loseInst.Process() {
			loseInst.End()
			return ErrLose
		}
	}

	if isRunAnim {
		if err := localanim.AnimMgrProcess(); err != nil {
			return fmt.Errorf("failed to handle animation: %w", err)
		}
	}

	battleCount++
	return nil
}

func Draw() {
	field.Draw()
	localanim.ObjAnimMgrDraw()
	localanim.AnimMgrDraw()

	drawEnemyNames()
	field.DrawBlackout()

	switch battleState {
	case stateOpening:
		if openingInst != nil {
			openingInst.Draw()
		}
	case stateChipSelect:
		chipsel.Draw()
	case stateBeforeMain:
		if b4mainInst != nil {
			b4mainInst.Draw()
		}
	case stateMain:
	case stateResultWin:
		win.Draw()
	case stateResultLose:
		if loseInst != nil {
			loseInst.Draw()
		}
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
	battleCount = 0
}

func drawEnemyNames() {
	for i, e := range enemyList {
		name := enemy.GetName(e.CharID)
		ofs := dxlib.GetDrawStringWidth(name, len(name))
		draw.String(config.ScreenSize.X-ofs-5, i*20+10, 0xffffff, "%s", name)
	}
}
