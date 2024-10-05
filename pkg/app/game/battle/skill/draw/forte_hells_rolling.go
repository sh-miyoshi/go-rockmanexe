package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawForteHellsRolling struct {
}

func (p *DrawForteHellsRolling) Draw(prev, current, next point.Point, count int, nextStepCount int, isFlip bool) {
	ofsx := battlecommon.GetOffset(next.X, current.X, prev.X, count%nextStepCount, nextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(next.Y, current.Y, prev.Y, count%nextStepCount, nextStepCount, battlecommon.PanelSize.Y)
	n := 0
	ydiff := next.Y - current.Y
	if ydiff != 0 {
		n = 2
	}
	opt := dxlib.DrawRotaGraphOption{}
	if isFlip {
		flag := int32(dxlib.TRUE)
		opt.ReverseXFlag = &flag
		ofsx -= 20
	} else {
		ofsx += 20
	}

	viewPos := battlecommon.ViewPos(current)
	dxlib.DrawRotaGraph(viewPos.X+ofsx, viewPos.Y+ofsy, 1.0, 0.0, images[imageTypeForteHellsRolling][n], true, opt)
}
