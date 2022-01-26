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
	absPlayerPosX float64
	absPlayerPosY float64
)

func Init() error {
	// TODO 本来ならplayerInfoから取得するが実装中なのでここでセットする
	var err error
	mapInfo, err = mapinfo.Load(mapinfo.IDTest)
	if err != nil {
		return fmt.Errorf("failed to load map info: %w", err)
	}
	absPlayerPosX = 300
	absPlayerPosY = 200

	collision.SetWalls(mapInfo.CollisionWalls)

	return nil
}

func End() {
	dxlib.DeleteGraph(mapInfo.Image)
}

func Draw() {
	var player, window common.Point
	getViewPos(&player, &window)

	dxlib.DrawRectGraph(0, 0, window.X, window.Y, mapInfo.Size.X, mapInfo.Size.Y, mapInfo.Image, true)

	if config.Get().Debug.ShowDebugData {
		// show debug data
		const color = 0xff0000
		dxlib.DrawCircle(player.X, player.Y, common.MapPlayerHitRange, color, true)
		for _, w := range mapInfo.CollisionWalls {
			cx := window.X
			cy := window.Y
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

	nextX, nextY := collision.NextPos(absPlayerPosX, absPlayerPosY, goVec)
	if nextX >= 0 && nextX < float64(mapInfo.Size.X) && nextY >= 0 && nextY < float64(mapInfo.Size.Y) {
		absPlayerPosX = nextX
		absPlayerPosY = nextY
	}

	return nil
}

func getViewPos(player, window *common.Point) {
	if absPlayerPosX < float64(common.ScreenSize.X/2) {
		player.X = int(absPlayerPosX)
		window.X = 0
	} else {
		// TODO 逆端
		player.X = common.ScreenSize.X / 2
		window.X = int(absPlayerPosX) - common.ScreenSize.X/2
	}

	if absPlayerPosY < float64(common.ScreenSize.Y/2) {
		player.Y = int(absPlayerPosY)
		window.Y = 0
	} else {
		// TODO 逆端
		player.Y = common.ScreenSize.Y / 2
		window.Y = int(absPlayerPosY) - common.ScreenSize.Y/2
	}
}
