package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawForteHellsRolling struct {
}

func (p *DrawForteHellsRolling) Draw(pos point.Point, count int, nextStepCount int) {
	n := count / 4 % len(imgForteHellsRolling)
	next := pos.X - 1
	prev := pos.X + 1
	ofsx := battlecommon.GetOffset(next, pos.X, prev, count%nextStepCount, nextStepCount, battlecommon.PanelSize.X) - battlecommon.PanelSize.X/2
	viewPos := battlecommon.ViewPos(pos)
	dxlib.DrawRotaGraph(viewPos.X+ofsx, viewPos.Y, 1.0, 0.0, imgForteHellsRolling[n], true)
}
