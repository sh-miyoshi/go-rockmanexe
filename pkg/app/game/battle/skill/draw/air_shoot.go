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

	// WIP imgNo
	dxlib.DrawRotaGraph(view.X+40, view.Y-15, 1, 0, images[imageTypeAirShootBody][0], true)
}
