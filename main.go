package main

import (
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
	flag.StringVar(&confFile, "config", "data/config.yaml", "file path of config")
	flag.Parse()

	if err := config.Init(confFile); err != nil {
		msg := fmt.Sprintf("failed to init config: %v", err)
		panic(msg)
	}

	rand.Seed(time.Now().Unix())
	dxlib.Init("data/DxLib.dll")

	if config.Get().Debug.Enabled {
		common.ImagePath = "data/private/images/"
		common.SoundPath = "data/private/sounds/"
		dxlib.SetOutApplicationLogValidFlag(dxlib.TRUE)
	} else {
		dxlib.SetOutApplicationLogValidFlag(dxlib.FALSE)
	}
	logger.InitLogger(config.Get().Debug.Enabled, config.Get().Log.FileName)

	fname := "data/font.ttf"
	if res := dxlib.AddFontFile(fname); res == nil {
		logger.Error("Failed to load font data %s", fname)
		return
	}

	dxlib.ChangeWindowMode(dxlib.TRUE)
	dxlib.SetGraphMode(common.ScreenX, common.ScreenY)

	dxlib.DxLib_Init()
	dxlib.SetDrawScreen(dxlib.DX_SCREEN_BACK)

	inputs.InitByDefault()
	if err := chip.Init("data/chipList.yaml"); err != nil {
		logger.Error("Failed to init chip data: %+v", err)
		return
	}
	if err := draw.Init(); err != nil {
		logger.Error("Failed to init drawing data: %+v", err)
		return
	}
	if err := sound.Init(); err != nil {
		logger.Error("Failed to init sound data: %+v", err)
		return
	}

	count := 0

MAIN:
	for dxlib.ScreenFlip() == 0 && dxlib.ProcessMessage() == 0 && dxlib.ClearDrawScreen() == 0 {
		inputs.KeyStateUpdate()
		if err := game.Process(); err != nil {
			logger.Error("Failed to play game: %+v", err)
			// TODO show to user
			break MAIN
		}
		game.Draw()

		if dxlib.CheckHitKey(dxlib.KEY_INPUT_ESCAPE) == 1 {
			break MAIN
		}
		count++

	}

	dxlib.DxLib_End()
}
