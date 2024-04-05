package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayCirkillShot = 4
)

type DrawCirkillShot struct {
}

func (p *DrawCirkillShot) Draw(prevPos, currentPos, nextPos point.Point, count int, nextStepCount int) {
	view := battlecommon.ViewPos(currentPos)
	n := (count / delayCirkillShot) % len(imgCirkillShot)

	cnt := count % nextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(nextPos.X, currentPos.X, prevPos.X, cnt, nextStepCount, battlecommon.PanelSize.X)
	xflip := int32(dxlib.TRUE)
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y, 1, 0, imgCirkillShot[n], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
}
