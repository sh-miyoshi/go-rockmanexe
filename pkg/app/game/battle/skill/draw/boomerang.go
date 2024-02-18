package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayBoomerang = 8
)

type DrawBoomerang struct {
}

func (p *DrawBoomerang) Draw(prevPos, currentPos, nextPos point.Point, count int, nextStepCount int) {
	view := battlecommon.ViewPos(currentPos)
	n := (count / delayBoomerang) % len(imgBoomerang)

	cnt := count % nextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(nextPos.X, currentPos.X, prevPos.X, cnt, nextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(nextPos.Y, currentPos.Y, prevPos.Y, cnt, nextStepCount, battlecommon.PanelSize.Y)
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+25+ofsy, 1, 0, imgBoomerang[n], true)
}
