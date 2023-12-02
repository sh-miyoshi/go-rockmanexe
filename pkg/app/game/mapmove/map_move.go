package mapmove

import (
	"errors"
	"fmt"
	"math"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/event"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/mapmove/collision"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/mapmove/scenario"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/vector"
)

var (
	ErrGoBattle = errors.New("go to battle")
	ErrGoMenu   = errors.New("go to menu")
	ErrGoEvent  = errors.New("go to event")

	mapInfo       mapinfo.MapInfo
	absPlayerPosX float64
	absPlayerPosY float64

	rockmanMoveImages      [5][]int
	rockmanMoveStandImages []int
	netMoveImages          [5][]int
	netMoveStandImages     []int

	playerMoveDirect int
	playerMoveCount  int

	initFlag bool = false
)

func Init() error {
	if initFlag {
		End()
		initFlag = false
	}

	// Load rockman image
	tmp := make([]int, 40)
	fname := common.ImagePath + "map/rockman_overworld_move.png"
	if res := dxlib.LoadDivGraph(fname, 30, 6, 5, 64, 64, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 6; j++ {
			rockmanMoveImages[i] = append(rockmanMoveImages[i], tmp[i*6+j])
		}
	}
	rockmanMoveStandImages = make([]int, 5)
	fname = common.ImagePath + "map/rockman_overworld_stand.png"
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 64, 64, rockmanMoveStandImages); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	// Load net image
	fname = common.ImagePath + "map/net_overworld_move.png"
	if res := dxlib.LoadDivGraph(fname, 40, 8, 5, 64, 64, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		for j := 0; j < 8; j++ {
			netMoveImages[i] = append(netMoveImages[i], tmp[i*8+j])
		}
	}
	netMoveStandImages = make([]int, 5)
	fname = common.ImagePath + "map/net_overworld_stand.png"
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 64, 64, netMoveStandImages); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	playerMoveDirect = common.DirectDown
	initFlag = true
	return nil
}

func End() {
	for i := 0; i < len(rockmanMoveImages); i++ {
		for _, img := range rockmanMoveImages[i] {
			dxlib.DeleteGraph(img)
		}
		rockmanMoveImages[i] = []int{}
	}
	for _, img := range rockmanMoveStandImages {
		dxlib.DeleteGraph(img)
	}
	rockmanMoveStandImages = []int{}

	for i := 0; i < len(netMoveImages); i++ {
		for _, img := range netMoveImages[i] {
			dxlib.DeleteGraph(img)
		}
		netMoveImages[i] = []int{}
	}
	for _, img := range netMoveStandImages {
		dxlib.DeleteGraph(img)
	}
	netMoveStandImages = []int{}
}

func MapChange(mapID int, pos point.Point) error {
	var err error
	mapInfo, err = mapinfo.Load(mapID)
	if err != nil {
		return fmt.Errorf("failed to load map info: %w", err)
	}
	absPlayerPosX = float64(pos.X)
	absPlayerPosY = float64(pos.Y)

	collision.SetWalls(mapInfo.CollisionWalls)
	collision.SetEvents(mapInfo.Events)
	logger.Info("change map to %d with %s", mapID, pos.String())
	return nil
}

func Draw() {
	var player, window point.Point
	getViewPos(&player, &window)

	dxlib.DrawRectGraph(0, 0, window.X, window.Y, mapInfo.Size.X, mapInfo.Size.Y, mapInfo.Image, true)

	drawPlayer(player)

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
	// デバッグ機能
	if inputs.CheckKey(inputs.KeyLButton) == 1 {
		// リロードの場合はyamlの情報も含めて再取得する
		mapinfo.Init(common.MapInfoFilePath)
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
		loadScenarioData(mapInfo.ID, e.No)
		return ErrGoEvent
	}

	if nextX >= 0 && nextX < float64(mapInfo.Size.X) && nextY >= 0 && nextY < float64(mapInfo.Size.Y) {
		absPlayerPosX = nextX
		absPlayerPosY = nextY
	}

	if nextDirect != 0 {
		// TODO(is_enemy_encounterがtrueなら敵遭遇処理)
		playerMoveCount++
		playerMoveDirect = nextDirect
	} else {
		playerMoveCount = 0
	}

	return nil
}

func getViewPos(player, window *point.Point) {
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

func drawPlayer(pos point.Point) {
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
		img := netMoveStandImages[typ]
		if mapInfo.IsCyberWorld {
			img = rockmanMoveStandImages[typ]
		}
		dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, img, true, dxopts)
	} else {
		n := (playerMoveCount / 4) % len(netMoveImages[typ])
		img := netMoveImages[typ][n]
		if mapInfo.IsCyberWorld {
			n = (playerMoveCount / 4) % len(rockmanMoveImages[typ])
			img = rockmanMoveImages[typ][n]
		}
		dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, img, true, dxopts)
	}
}

func loadScenarioData(mapType int, eventNo int) {
	logger.Info("load scenario for map %d, event %d", mapType, eventNo)

	switch mapType {
	case mapinfo.ID_犬小屋:
		event.SetScenarios(scenario.Scenario_犬小屋[eventNo])
	case mapinfo.ID_秋原町:
		event.SetScenarios(scenario.Scenario_秋原町[eventNo])
	default:
		common.SetError(fmt.Sprintf("no scenario data for map type %d", mapType))
	}
}
