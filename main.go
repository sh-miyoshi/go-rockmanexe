package main

import (
	"runtime"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/game"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	dxlib.Init("DxLib.dll")

	dxlib.ChangeWindowMode(dxlib.TRUE)
	dxlib.SetGraphMode(common.ScreenX, common.ScreenY)
	dxlib.SetOutApplicationLogValidFlag(dxlib.TRUE)

	dxlib.DxLib_Init()
	dxlib.SetDrawScreen(dxlib.DX_SCREEN_BACK)

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
