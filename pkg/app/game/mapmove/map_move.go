package mapmove

import (
	"errors"
	"fmt"
	"math"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/event"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/mapmove/collision"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/vector"
)

var (
	ErrGoBattle = errors.New("go to battle")
	ErrGoMenu   = errors.New("go to menu")
	ErrGoEvent  = errors.New("go to event")

	mapInfo       *mapinfo.MapInfo
	absPlayerPosX float64
	absPlayerPosY float64

	playerMoveImages      [5][]int
	playerMoveStandImages []int
	playerMoveDirect      int
	playerMoveCount       int
)

func Init() error {
	// TODO 本来ならplayerInfoから取得するが実装中なのでここでセットする
	var err error
	mapInfo, err = mapinfo.Load(mapinfo.ID_犬小屋)
	if err != nil {
		return fmt.Errorf("failed to load map info: %w", err)
	}
	absPlayerPosX = 300
	absPlayerPosY = 200

	collision.SetWalls(mapInfo.CollisionWalls)
	collision.SetEvents(mapInfo.Events)

	// Load player image
	tmp := make([]int, 30)
	fname := common.ImagePath + "map/rockman_overworld_move.png"
	if res := dxlib.LoadDivGraph(fname, 30, 6, 5, 64, 64, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			playerMoveImages[i] = append(playerMoveImages[i], tmp[i*6+j])
		}
	}
	playerMoveStandImages = make([]int, 5)
	fname = common.ImagePath + "map/rockman_overworld_stand.png"
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 64, 64, playerMoveStandImages); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	playerMoveDirect = common.DirectDown

	return nil
}

func End() {
	dxlib.DeleteGraph(mapInfo.Image)
	for i := 0; i < 5; i++ {
		for _, img := range playerMoveImages[i] {
			dxlib.DeleteGraph(img)
		}
		playerMoveImages[i] = []int{}
	}
}

func Draw() {
	var player, window common.Point
	getViewPos(&player, &window)

	dxlib.DrawRectGraph(0, 0, window.X, window.Y, mapInfo.Size.X, mapInfo.Size.Y, mapInfo.Image, true)

	drawRockman(player)

	if config.Get().Debug.ShowDebugData {
		// show debug data
		const color = 0xff0000
		dxlib.DrawCircle(player.X, player.Y, common.MapPlayerHitRange, color, true)
		for _, w := range mapInfo.CollisionWalls {
			cx := window.X
			cy := window.Y
			dxlib.DrawLine(w.X1-cx, w.Y1-cy, w.X2-cx, w.Y2-cy, color)
		}

		for _, e := range mapInfo.Events {
			cx := window.X
			cy := window.Y
			dxlib.DrawCircle(e.X-cx, e.Y-cy, e.R, color, false)
		}

		dxlib.DrawFormatString(0, 0, color, "Window: (%d, %d)", window.X, window.Y)
		dxlib.DrawFormatString(0, 20, color, "Player: (%d, %d)", player.X, player.Y)
		dxlib.DrawFormatString(0, 40, color, "ABS: (%.2f, %.2f)", absPlayerPosX, absPlayerPosY)
		dxlib.DrawFormatString(0, 60, color, "Reload: L-btn")
	}
}

func Process() error {
	if inputs.CheckKey(inputs.KeyLButton) == 1 {
		End()
		Init()
		return nil
	}

	goVec := vector.Vector{}
	nextDirect := 0
	if inputs.CheckKey(inputs.KeyRight) != 0 {
		goVec.X += 4
		nextDirect |= common.DirectRight
	} else if inputs.CheckKey(inputs.KeyLeft) != 0 {
		goVec.X -= 4
		nextDirect |= common.DirectLeft
	}

	if inputs.CheckKey(inputs.KeyDown) != 0 {
		goVec.Y += 4
		nextDirect |= common.DirectDown
	} else if inputs.CheckKey(inputs.KeyUp) != 0 {
		goVec.Y -= 4
		nextDirect |= common.DirectUp
	}

	// 斜め移動の場合は速度を調整
	if goVec.X != 0 && goVec.Y != 0 {
		goVec.X /= math.Sqrt(2)
		goVec.Y /= math.Sqrt(2)
	}

	nextX, nextY := collision.NextPos(absPlayerPosX, absPlayerPosY, goVec)
	if e := collision.GetEvent(nextX, nextY); e != nil {
		// Hit to Event
		event.Set(e.Type, e.Args)
		return ErrGoEvent
	}
	// TODO(hit events, or object)

	if nextX >= 0 && nextX < float64(mapInfo.Size.X) && nextY >= 0 && nextY < float64(mapInfo.Size.Y) {
		absPlayerPosX = nextX
		absPlayerPosY = nextY
	}

	if nextDirect != 0 {
		playerMoveCount++
		playerMoveDirect = nextDirect
	} else {
		playerMoveCount = 0
	}

	return nil
}

func getViewPos(player, window *common.Point) {
	hsX := common.ScreenSize.X / 2
	hsY := common.ScreenSize.Y / 2

	if absPlayerPosX < float64(hsX) {
		player.X = int(absPlayerPosX)
		window.X = 0
	} else {
		s := mapInfo.Size.X - hsX
		if absPlayerPosX > float64(s) {
			window.X = mapInfo.Size.X - common.ScreenSize.X
			player.X = hsX + int(absPlayerPosX) - s
		} else {
			player.X = hsX
			window.X = int(absPlayerPosX) - hsX
		}
	}

	if absPlayerPosY < float64(hsY) {
		player.Y = int(absPlayerPosY)
		window.Y = 0
	} else {
		s := mapInfo.Size.Y - hsY
		if absPlayerPosY > float64(s) {
			window.Y = mapInfo.Size.Y - common.ScreenSize.Y
			player.Y = hsY + int(absPlayerPosY) - s
		} else {
			player.Y = hsY
			window.Y = int(absPlayerPosY) - hsY
		}
	}
}

func drawRockman(pos common.Point) {
	dxopts := dxlib.DrawRotaGraphOption{}
	rlFlag := false
	if playerMoveDirect&common.DirectLeft != 0 {
		t := int32(dxlib.TRUE)
		dxopts.ReverseXFlag = &t
		rlFlag = true
	} else if playerMoveDirect&common.DirectRight != 0 {
		rlFlag = true
	}

	typ := 0
	if playerMoveDirect&common.DirectUp != 0 {
		typ = 0
		if rlFlag {
			typ = 1
		}
	} else if playerMoveDirect&common.DirectDown != 0 {
		typ = 4
		if rlFlag {
			typ = 3
		}
	} else {
		typ = 2
	}

	if playerMoveCount == 0 {
		dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, playerMoveStandImages[typ], true, dxopts)
	} else {
		n := (playerMoveCount / 4) % 6
		dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, playerMoveImages[typ][n], true, dxopts)
	}
}
