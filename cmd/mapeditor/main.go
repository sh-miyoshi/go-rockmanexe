package main

import (
	"fmt"
	"os"
	"runtime"

	origindxlib "github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"gopkg.in/yaml.v2"
)

var (
	window  point.Point
	mapInfo mapinfo.MapInfo
)

func init() {
	runtime.LockOSThread()
}

func main() {
	origindxlib.Init("../../data/DxLib.dll")

	origindxlib.ChangeWindowMode(origindxlib.TRUE)
	origindxlib.SetGraphMode(960, 640)
	origindxlib.SetOutApplicationLogValidFlag(origindxlib.TRUE)

	origindxlib.DxLib_Init()
	origindxlib.SetDrawScreen(origindxlib.DX_SCREEN_BACK)

	image := dxlib.LoadGraph("base.png")
	if image == -1 {
		fmt.Println("Failed to load base.png")
		return
	}

	var imgSize point.Point
	dxlib.GetGraphSize(image, &imgSize.X, &imgSize.Y)

	readWalls()

MAIN:
	for origindxlib.ScreenFlip() == 0 && origindxlib.ProcessMessage() == 0 && origindxlib.ClearDrawScreen() == 0 {
		if origindxlib.CheckHitKey(origindxlib.KEY_INPUT_ESCAPE) == 1 {
			break MAIN
		}

		dxlib.DrawRectGraph(0, 0, window.X, window.Y, imgSize.X, imgSize.Y, image, true)
		for _, w := range mapInfo.CollisionWalls {
			dxlib.DrawLine(w.X1-window.X, w.Y1-window.Y, w.X2-window.X, w.Y2-window.Y, 0xff0000)
		}

		var mouseX, mouseY int32
		origindxlib.GetMousePoint(&mouseX, &mouseY)
		origindxlib.DrawFormatString(5, 615, 0xffffff, "( %d, %d )", mouseX+int32(window.X), mouseY+int32(window.Y))

		spd := 10
		if inputs.CheckKey(inputs.KeyLeft) > 0 && window.X > spd {
			window.X -= spd
		}
		if inputs.CheckKey(inputs.KeyRight) > 0 && window.X < imgSize.X-spd {
			window.X += spd
		}
		if inputs.CheckKey(inputs.KeyUp) > 0 && window.Y > spd {
			window.Y -= spd
		}
		if inputs.CheckKey(inputs.KeyDown) > 0 && window.Y < imgSize.Y-spd {
			window.Y += spd
		}

		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			readWalls()
		}
	}

	origindxlib.DxLib_End()
}

func readWalls() {
	fp, err := os.Open("info.yaml")
	if err != nil {
		common.SetError(err.Error())
	}
	defer fp.Close()

	if err := yaml.NewDecoder(fp).Decode(&mapInfo); err != nil {
		common.SetError(err.Error())
	}
}
