package main

import (
	"os"
	"runtime"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"gopkg.in/yaml.v2"
)

var (
	playerID string
)

const (
	baseDir = "../../"
)

type appConfig struct {
	Fps int64 `yaml:"fps"`
}

func init() {
	runtime.LockOSThread()
}

func main() {
	fps.FPS = 60
	fpsMgr := fps.Fps{}

	dxlib.Init(baseDir + config.DxlibDLLFilePath)

	dxlib.SetDoubleStartValidFlag(dxlib.TRUE)
	if config.Get().Debug.RunAlways {
		dxlib.SetAlwaysRunFlag(dxlib.TRUE)
	}

	dxlib.SetOutApplicationLogValidFlag(dxlib.TRUE)
	logger.InitLogger(true, "application.log")

	if res := dxlib.AddFontFile(baseDir + config.FontFilePath); res == nil {
		logger.Error("Failed to load font data %s", config.FontFilePath)
		return
	}

	dxlib.ChangeWindowMode(dxlib.TRUE)
	dxlib.SetWindowSizeChangeEnableFlag(dxlib.FALSE, dxlib.FALSE)
	dxlib.SetGraphMode(config.MaxScreenSize.X, config.MaxScreenSize.Y)
	dxlib.SetWindowSize(int32(config.ScreenSize.X), int32(config.ScreenSize.Y))

	dxlib.DxLib_Init()
	dxlib.SetDrawScreen(dxlib.DX_SCREEN_BACK)

	count := 0

	if err := appInit(); err != nil {
		logger.Error("Failed to init application: %+v", err)
		system.SetError("ゲーム初期化時")
	}

	logger.Info("Successfully init application.")

	pos := point.Point{X: 1, Y: 1}
	var act battleplayer.BattlePlayerAct
	act.Init(&pos)

	chipID := chip.IDCannon

MAIN:
	for system.Error() == nil && dxlib.ScreenFlip() == 0 && dxlib.ProcessMessage() == 0 && dxlib.ClearDrawScreen() == 0 {
		inputs.KeyStateUpdate()

		// メインロジック
		background.Update()
		field.Update()
		background.Draw()
		field.Draw()
		localanim.ObjAnimMgrProcess(false, field.IsBlackout())
		playerDraw(pos, act)
		localanim.AnimMgrProcess()
		// localanim.ObjAnimMgrDraw()
		localanim.AnimMgrDraw()

		if !act.Update() && inputs.CheckKey(inputs.KeyEnter) == 1 {
			cfg := loadConfig()
			fps.FPS = cfg.Fps

			c := chip.Get(chipID)
			if c.PlayerAct != -1 {
				act.SetAnim(c.PlayerAct, c.KeepCount)
			}
			sid := skillcore.GetIDByChipID(c.ID)
			act.SetSkill(sid, skillcore.Argument{
				OwnerID:    playerID,
				Power:      c.Power,
				TargetType: 0,
			})
		}

		if dxlib.CheckHitKey(dxlib.KEY_INPUT_ESCAPE) == 1 {
			logger.Info("Game end by escape command")
			break MAIN
		}
		count++

		fpsMgr.Wait()
	}

	if err := system.Error(); err != nil {
		sound.BGMStop()
		dxlib.ClearDrawScreen()
		dxlib.DrawFormatString(10, 10, 0xff0000, "%sに回復不可能なエラーが発生しました。", err.Error())
		dxlib.DrawFormatString(10, 40, 0xff0000, "詳細はログを参照してください。")
		dxlib.ScreenFlip()
		dxlib.WaitKey()
	}

	dxlib.DxLib_End()
}

func appInit() error {
	config.ImagePath = baseDir + config.ImagePath

	if err := inputs.Init(config.Get().Input.Type); err != nil {
		return errors.Wrap(err, "inputs init failed")
	}
	if err := chip.Init(baseDir + config.ChipFilePath); err != nil {
		return errors.Wrap(err, "chip init failed")
	}
	if err := draw.Init(); err != nil {
		return errors.Wrap(err, "drawing data init failed")
	}

	if err := field.Init(); err != nil {
		return errors.Wrap(err, "battle field init failed")
	}

	if err := skill.Init(); err != nil {
		return errors.Wrap(err, "skill init failed")
	}

	playerInst, _ := battleplayer.New(&player.Player{
		HP:         100,
		ShotPower:  1,
		ChargeTime: 120,
	})
	playerID = playerInst.ID
	localanim.ObjAnimNew(playerInst)

	return nil
}

func playerDraw(pos point.Point, act battleplayer.BattlePlayerAct) {
	view := battlecommon.ViewPos(pos)
	img := act.GetImage()
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, img, true)
}

func loadConfig() appConfig {
	fp, _ := os.Open("config.yaml")
	defer fp.Close()

	var conf appConfig

	yaml.NewDecoder(fp).Decode(&conf)
	return conf
}
