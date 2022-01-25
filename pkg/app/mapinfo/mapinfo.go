package mapinfo

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type Wall struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

type MapInfo struct {
	Image          int
	Size           common.Point
	CollisionWalls []Wall
}

const (
	IDTest int = iota
)

func Load(id int) (*MapInfo, error) {
	basePath := common.ImagePath + "map/"
	var fname string
	collisionWalls := []Wall{}

	switch id {
	case IDTest:
		fname = basePath + "test.png"
		collisionWalls = []Wall{
			{X1: 15, Y1: 327, X2: 649, Y2: 9},
			{X1: 649, Y1: 9, X2: 1285, Y2: 326},
			{X1: 1285, Y1: 326, X2: 652, Y2: 644},
			{X1: 652, Y1: 644, X2: 15, Y2: 327},
		}
	default:
		return nil, fmt.Errorf("no such map: %d", id)
	}
	res := &MapInfo{
		Image:          dxlib.LoadGraph(fname),
		CollisionWalls: collisionWalls,
	}
	if res.Image == -1 {
		return nil, fmt.Errorf("failed to load image: %s", fname)
	}
	dxlib.GetGraphSize(res.Image, &res.Size.X, &res.Size.Y)

	return res, nil
}
