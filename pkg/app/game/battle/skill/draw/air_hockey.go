package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawAirHockey struct {
}

func (p *DrawAirHockey) Draw(prevPos, currentPos, nextPos point.Point, count int, nextStepCount int) {
	view := battlecommon.ViewPos(currentPos)

	cnt := count % nextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(nextPos.X, currentPos.X, prevPos.X, cnt, nextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(nextPos.Y, currentPos.Y, prevPos.Y, cnt, nextStepCount, battlecommon.PanelSize.Y)
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+30+ofsy, 1, 0, images[imageTypeAirHockey][0], true)
}
