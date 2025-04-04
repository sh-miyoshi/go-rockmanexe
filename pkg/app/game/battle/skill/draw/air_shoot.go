package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawAirShoot struct {
}

func (p *DrawAirShoot) Draw(objPos point.Point, count int) {
	view := battlecommon.ViewPos(objPos)

	imageNo := 0
	if count >= 8 {
		imageNo = 1
	}
	dxlib.DrawRotaGraph(view.X+40, view.Y-15, 1, 0, images[imageTypeAirShootBody][imageNo], true)

	imageNo = count / 3
	if imageNo < len(images[imageTypeAirShootAtk]) {
		dxlib.DrawRotaGraph(view.X+80, view.Y-15, 1, 0, images[imageTypeAirShootAtk][imageNo], true)
	}
}
