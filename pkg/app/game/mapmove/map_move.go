package mapmove

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/mapmove/collision"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/vector"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

var (
	ErrGoBattle = errors.New("go to battle")
	ErrGoMenu   = errors.New("go to menu")

	mapInfo       *mapinfo.MapInfo
	currentWindow common.Point
	playerPos     common.Point
)

func Init() error {
	// TODO 本来ならplayerInfoから取得するが実装中なのでここでセットする
	var err error
	mapInfo, err = mapinfo.Load(mapinfo.IDTest)
	if err != nil {
		return fmt.Errorf("failed to load map info: %w", err)
	}
	currentWindow.X = 100
	currentWindow.Y = 100
	playerPos.X = common.ScreenSize.X / 2
	playerPos.Y = common.ScreenSize.Y / 2

	collision.SetWalls(mapInfo.CollisionWalls)

	return nil
}

func End() {
	dxlib.DeleteGraph(mapInfo.Image)
}

func Draw() {
	dxlib.DrawRectGraph(0, 0, currentWindow.X, currentWindow.Y, mapInfo.Size.X, mapInfo.Size.Y, mapInfo.Image, true)

	if config.Get().Debug.ShowDebugData {
		// show debug data
		const color = 0xff0000
		dxlib.DrawCircle(playerPos.X, playerPos.Y, common.MapPlayerHitRange, color, true)
		for _, w := range mapInfo.CollisionWalls {
			cx := currentWindow.X
			cy := currentWindow.Y
			dxlib.DrawLine(w.X1-cx, w.Y1-cy, w.X2-cx, w.Y2-cy, color)
		}
	}
}

func Process() error {
	goVec := vector.Vector{}
	if inputs.CheckKey(inputs.KeyRight) != 0 {
		goVec.X += 4
	}
	if inputs.CheckKey(inputs.KeyLeft) != 0 {
		goVec.X -= 4
	}
	if inputs.CheckKey(inputs.KeyDown) != 0 {
		goVec.Y += 4
	}
	if inputs.CheckKey(inputs.KeyUp) != 0 {
		goVec.Y -= 4
	}

	next := collision.NextPos(playerPos, goVec)
	playerPos = next

	// TODO
	//current := playerPos.Add(currentWindow)
	// if (goVec.X > 0 && p.X <= common.ScreenSize.X/2) || (goVec.X < 0 && p.X >= common.ScreenSize.X/2) {
	// 	playerPos.X = p.X
	// } else if w.X >= 0 && w.X <= mapInfo.Size.X {
	// 	currentWindow.X = w.X
	// } else if p.X >= 0 && p.X <= common.ScreenSize.X {
	// 	playerPos.X = p.X
	// }

	return nil
}
