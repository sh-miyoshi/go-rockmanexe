package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

var (
	version = ""
	encKey  = ""
)

func init() {
	runtime.LockOSThread()
}

func main() {
	fpsMgr := fps.Fps{TargetFPS: 60}

	var confFile string
	flag.StringVar(&confFile, "config", common.DefaultConfigFile, "file path of config")
	flag.Parse()

	if err := config.Init(confFile); err != nil {
		msg := fmt.Sprintf("failed to init config: %v", err)
		panic(msg)
	}

	rand.Seed(time.Now().Unix())
	dxlib.Init(common.DxlibDLLFilePath)

	if version != "" {
		common.ProgramVersion = version
	}
	if encKey != "" {
		common.EncryptKey = encKey
	}

	dxlib.SetDoubleStartValidFlag(dxlib.TRUE)
	if config.Get().Debug.RunAlways {
		dxlib.SetAlwaysRunFlag(dxlib.TRUE)
	}

	if config.Get().Debug.UsePrivateResource {
		common.ImagePath = "data/private/images/"
		common.SoundPath = "data/private/sounds/"
	}

	if config.Get().Log.DebugEnabled {
		dxlib.SetOutApplicationLogValidFlag(dxlib.TRUE)
	} else {
		dxlib.SetOutApplicationLogValidFlag(dxlib.FALSE)
	}
	logger.InitLogger(config.Get().Log.DebugEnabled, config.Get().Log.FileName)

	if res := dxlib.AddFontFile(common.FontFilePath); res == nil {
		logger.Error("Failed to load font data %s", common.FontFilePath)
		return
	}

	dxlib.ChangeWindowMode(dxlib.TRUE)
	dxlib.SetWindowSizeChangeEnableFlag(dxlib.FALSE, dxlib.FALSE)
	dxlib.SetGraphMode(int32(common.MaxScreenSize.X), int32(common.MaxScreenSize.Y))
	dxlib.SetWindowSize(int32(common.ScreenSize.X), int32(common.ScreenSize.Y))

	dxlib.DxLib_Init()
	dxlib.SetDrawScreen(dxlib.DX_SCREEN_BACK)

	logger.Info("Program version: %s", common.ProgramVersion)

	count := 0

	if err := appInit(); err != nil {
		logger.Error("Failed to init application: %+v", err)
		common.IrreversibleError = errors.New("ゲーム初期化時")
	}

	logger.Info("Successfully init application.")

	if config.Get().Debug.InitSleepSec > 0 {
		tm := time.Duration(config.Get().Debug.InitSleepSec) * time.Second
		time.Sleep(tm)
	}
MAIN:
	for common.IrreversibleError == nil && dxlib.ScreenFlip() == 0 && dxlib.ProcessMessage() == 0 && dxlib.ClearDrawScreen() == 0 {
		inputs.KeyStateUpdate()
		if err := game.Process(); err != nil {
			logger.Error("Failed to play game: %+v", err)
			common.IrreversibleError = errors.New("ゲームプレイ中")
			break MAIN
		}
		game.Draw()

		if dxlib.CheckHitKey(dxlib.KEY_INPUT_ESCAPE) == 1 {
			logger.Info("Game end by escape command")
			break MAIN
		}
		count++

		fpsMgr.Wait()
		if config.Get().Debug.ShowDebugData {
			dxlib.DrawFormatString(int32(common.ScreenSize.X-60), 10, 0xff0000, "[%.1f]", fpsMgr.Get())
		}
	}

	if common.IrreversibleError != nil {
		sound.BGMStop()
		dxlib.ClearDrawScreen()
		dxlib.DrawFormatString(10, 10, 0xff0000, "%sに回復不可能なエラーが発生しました。", common.IrreversibleError.Error())
		dxlib.DrawFormatString(10, 40, 0xff0000, "詳細はログを参照してください。")
		dxlib.ScreenFlip()
		dxlib.WaitKey()
	}

	dxlib.DxLib_End()
}

func appInit() error {
	if err := inputs.Init(config.Get().Input.Type); err != nil {
		return fmt.Errorf("inputs init failed: %w", err)
	}
	if err := chip.Init(common.ChipFilePath); err != nil {
		return fmt.Errorf("chip init failed: %w", err)
	}
	if err := draw.Init(); err != nil {
		return fmt.Errorf("drawing data init failed: %w", err)
	}
	if err := sound.Init(); err != nil {
		return fmt.Errorf("sound init failed: %w", err)
	}
	if err := mapinfo.Init(common.MapInfoFilePath); err != nil {
		return fmt.Errorf("map info init failed: %w", err)
	}

	return nil
}
