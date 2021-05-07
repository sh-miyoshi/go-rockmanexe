package main

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/game"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	var confFile string
	flag.StringVar(&confFile, "config", common.DefaultConfigFile, "file path of config")
	flag.Parse()

	if err := config.Init(confFile); err != nil {
		msg := fmt.Sprintf("failed to init config: %v", err)
		panic(msg)
	}

	rand.Seed(time.Now().Unix())
	dxlib.Init(common.DxlibDLLFilePath)

	if config.Get().Debug.Enabled {
		common.ImagePath = "data/private/images/"
		common.SoundPath = "data/private/sounds/"
		dxlib.SetOutApplicationLogValidFlag(dxlib.TRUE)
	} else {
		dxlib.SetOutApplicationLogValidFlag(dxlib.FALSE)
	}
	logger.InitLogger(config.Get().Debug.Enabled, config.Get().Log.FileName)

	if res := dxlib.AddFontFile(common.FontFilePath); res == nil {
		logger.Error("Failed to load font data %s", common.FontFilePath)
		return
	}

	dxlib.ChangeWindowMode(dxlib.TRUE)
	dxlib.SetGraphMode(common.ScreenX, common.ScreenY)

	dxlib.DxLib_Init()
	dxlib.SetDrawScreen(dxlib.DX_SCREEN_BACK)

	logger.Info("Program version: %s", common.ProgramVersion)

	count := 0
	var exitErr error

	if err := appInit(); err != nil {
		logger.Error("Failed to init application: %+v", err)
		exitErr = errors.New("ゲーム初期化時")
	}

	logger.Info("Successfully init application.")
MAIN:
	for exitErr == nil && dxlib.ScreenFlip() == 0 && dxlib.ProcessMessage() == 0 && dxlib.ClearDrawScreen() == 0 {
		inputs.KeyStateUpdate()
		if err := game.Process(); err != nil {
			logger.Error("Failed to play game: %+v", err)
			exitErr = errors.New("ゲームプレイ中")
			break MAIN
		}
		game.Draw()

		if dxlib.CheckHitKey(dxlib.KEY_INPUT_ESCAPE) == 1 {
			logger.Info("Game end by escape command")
			break MAIN
		}
		count++

	}

	if exitErr != nil {
		dxlib.ClearDrawScreen()
		dxlib.DrawFormatString(10, 10, 0xff0000, "%sに回復不可能なエラーが発生しました。", exitErr.Error())
		dxlib.DrawFormatString(10, 40, 0xff0000, "詳細はログを参照してください。")
		dxlib.ScreenFlip()
		dxlib.WaitKey()
	}

	dxlib.DxLib_End()
}

func appInit() error {
	inputs.InitByDefault()
	if err := chip.Init(common.ChipFilePath); err != nil {
		return fmt.Errorf("chip init failed: %w", err)
	}
	if err := draw.Init(); err != nil {
		return fmt.Errorf("drawing data init failed: %w", err)
	}
	if err := sound.Init(); err != nil {
		return fmt.Errorf("sound init failed: %w", err)
	}

	return nil
}
