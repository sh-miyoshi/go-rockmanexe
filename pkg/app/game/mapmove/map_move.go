package mapmove

import (
	"errors"
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/mapinfo"
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

	return nil
}

func End() {
	dxlib.DeleteGraph(mapInfo.Image)
}

func Draw() {
	dxlib.DrawRectGraph(0, 0, currentWindow.X, currentWindow.Y, mapInfo.Size.X, mapInfo.Size.Y, mapInfo.Image, true)
}

func Process() error {
	return nil
}
