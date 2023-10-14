package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayGarooBreath = 4
)

type DrawGarooBreath struct {
}

func (p *DrawGarooBreath) Draw(prevPos, currentPos, nextPos common.Point, count int) {
	view := battlecommon.ViewPos(currentPos)
	n := (count / delayGarooBreath) % len(imgGarooBreath)

	cnt := count % resources.SkillGarooBreathNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(nextPos.X, currentPos.X, prevPos.X, cnt, resources.SkillGarooBreathNextStepCount, battlecommon.PanelSize.X)
	ofsy := -15
	xflip := int32(dxlib.TRUE)
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, imgGarooBreath[n], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
}
