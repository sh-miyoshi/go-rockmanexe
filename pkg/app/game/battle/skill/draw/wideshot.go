package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawWideShot struct {
}

func (p *DrawWideShot) Draw(pos common.Point, count int, direct int, showBody bool, nextStepCount int, state int) {
	opt := dxlib.DrawRotaGraphOption{}
	ofs := 1
	if direct == common.DirectLeft {
		xflip := int32(dxlib.TRUE)
		opt.ReverseXFlag = &xflip
		ofs = -1
	}
	view := battlecommon.ViewPos(pos)

	switch state {
	case resources.SkillWideShotStateBegin:
		n := (count / resources.SkillWideShotDelay)

		if n < len(imgWideShotBody) && showBody {
			dxlib.DrawRotaGraph(view.X+40, view.Y-13, 1, 0, imgWideShotBody[n], true, opt)
		}
		if n >= len(imgWideShotBegin) {
			n = len(imgWideShotBegin) - 1
		}
		dxlib.DrawRotaGraph(view.X+62*ofs, view.Y+20, 1, 0, imgWideShotBegin[n], true, opt)
	case resources.SkillWideShotStateMove:
		n := (count / resources.SkillWideShotDelay) % len(imgWideShotMove)
		next := pos.X + 1
		prev := pos.X - 1
		if direct == common.DirectLeft {
			next, prev = prev, next
		}

		c := count % nextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, pos.X, prev, c, nextStepCount, battlecommon.PanelSize.X)
			dxlib.DrawRotaGraph(view.X+ofsx, view.Y+20, 1, 0, imgWideShotMove[n], true, opt)
		}
	}
}
