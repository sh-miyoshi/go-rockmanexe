package mapinfo

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type MapInfo struct {
	Image        int
	Size         common.Point
	CollisionPos []common.Point
}

const (
	IDTest int = iota
)

func Load(id int) (*MapInfo, error) {
	basePath := common.ImagePath + "map/"
	var fname string

	switch id {
	case IDTest:
		fname = basePath + "test.png"
	default:
		return nil, fmt.Errorf("no such map: %d", id)
	}
	res := &MapInfo{
		Image: dxlib.LoadGraph(fname),
	}
	if res.Image == -1 {
		return nil, fmt.Errorf("failed to load image: %s", fname)
	}
	dxlib.GetGraphSize(res.Image, &res.Size.X, &res.Size.Y)

	return res, nil
}
