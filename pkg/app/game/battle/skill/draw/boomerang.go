package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayBoomerang = 8
)

type DrawBoomerang struct {
}

func (p *DrawBoomerang) Draw(prevPos, currentPos, nextPos common.Point, count int) {
	view := battlecommon.ViewPos(currentPos)
	n := (count / delayBoomerang) % len(imgBoomerang)

	cnt := count % resources.SkillBoomerangNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(nextPos.X, currentPos.X, prevPos.X, cnt, resources.SkillBoomerangNextStepCount, battlecommon.PanelSize.X)
	ofsy := battlecommon.GetOffset(nextPos.Y, currentPos.Y, prevPos.Y, cnt, resources.SkillBoomerangNextStepCount, battlecommon.PanelSize.Y)
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+25+ofsy, 1, 0, imgBoomerang[n], true)
}
