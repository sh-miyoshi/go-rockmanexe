package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawForteHellsRolling struct {
}

func (p *DrawForteHellsRolling) Draw(prev, current, next point.Point, count int, nextStepCount int) {
	ofsx := battlecommon.GetOffset(next.X, current.X, prev.X, count%nextStepCount, nextStepCount, battlecommon.PanelSize.X) - battlecommon.PanelSize.X/2
	ofsy := battlecommon.GetOffset(next.Y, current.Y, prev.Y, count%nextStepCount, nextStepCount, battlecommon.PanelSize.Y)
	ydiff := next.Y - current.Y
	n := 0
	angle := 0.0
	if ydiff > 0 {
		n = 2
		angle = -45.0
	} else if ydiff < 0 {
		n = 2
		angle = 45.0
	}
	viewPos := battlecommon.ViewPos(current)
	dxlib.DrawRotaGraph(viewPos.X+ofsx, viewPos.Y+ofsy, 1.0, angle, imgForteHellsRolling[n], true)
}
