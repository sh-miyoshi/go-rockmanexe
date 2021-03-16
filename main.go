package main

import (
	"flag"
	"runtime"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/game"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run as debug mode")
	flag.Parse()

	if debug {
		common.ImagePath = "data/private/images/"
		logger.InitLogger(true, "")
	}

	dxlib.Init("DxLib.dll")

	fname := "data/font.ttf"
	if res := dxlib.AddFontFile(fname); res == nil {
		logger.Error("Failed to load font data %s", fname)
		return
	}

	dxlib.ChangeWindowMode(dxlib.TRUE)
	dxlib.SetGraphMode(common.ScreenX, common.ScreenY)
	dxlib.SetOutApplicationLogValidFlag(dxlib.TRUE)

	dxlib.DxLib_Init()
	dxlib.SetDrawScreen(dxlib.DX_SCREEN_BACK)

	inputs.InitByDefault()
	if err := chip.Init("data/chipList.yaml"); err != nil {
		logger.Error("Failed to init chip data: %v", err)
		return
	}
	if err := draw.Init(); err != nil {
		logger.Error("Failed to init drawing data: %v", err)
		return
	}

	count := 0

MAIN:
	for dxlib.ScreenFlip() == 0 && dxlib.ProcessMessage() == 0 && dxlib.ClearDrawScreen() == 0 {
		inputs.KeyStateUpdate()
		if err := game.Process(); err != nil {
			logger.Error("Failed to play game: %v", err)
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
