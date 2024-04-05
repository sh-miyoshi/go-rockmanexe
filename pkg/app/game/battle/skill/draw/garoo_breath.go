package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayGarooBreath = 4
)

type DrawGarooBreath struct {
}

func (p *DrawGarooBreath) Draw(prevPos, currentPos, nextPos point.Point, count int, nextStepCount int) {
	view := battlecommon.ViewPos(currentPos)
	n := (count / delayGarooBreath) % len(imgGarooBreath)

	cnt := count % nextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(nextPos.X, currentPos.X, prevPos.X, cnt, nextStepCount, battlecommon.PanelSize.X)
	ofsy := -15
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, imgGarooBreath[n], true, dxlib.OptXReverse(true))
}
