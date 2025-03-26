package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawAirHockey struct {
}

func (p *DrawAirHockey) Draw(prevPos, currentPos, nextPos point.Point, count int, nextStepCount int) {
	// WIP
	view := battlecommon.ViewPos(currentPos)
	dxlib.DrawRotaGraph(view.X, view.Y+25, 1, 0, images[imageTypeAirHockey][0], true)
}
