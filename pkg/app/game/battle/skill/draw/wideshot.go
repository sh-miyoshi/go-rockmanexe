package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayWideShot = 4
)

type DrawWideShot struct {
}

func (p *DrawWideShot) Draw(pos point.Point, count int, direct int, showBody bool, nextStepCount int, state int) {
	opt := dxlib.OptXReverse(direct == config.DirectLeft)
	view := battlecommon.ViewPos(pos)

	switch state {
	case resources.SkillWideShotStateBegin:
		n := (count / delayWideShot)

		if n < len(images[imageTypeWideShotBody]) && showBody {
			dxlib.DrawRotaGraph(view.X+math.ReverseIf(40, direct == config.DirectLeft), view.Y-13, 1, 0, images[imageTypeWideShotBody][n], true, opt)
		}
		if n >= len(images[imageTypeWideShotBegin]) {
			n = len(images[imageTypeWideShotBegin]) - 1
		}
		dxlib.DrawRotaGraph(view.X+math.ReverseIf(62, direct == config.DirectLeft), view.Y+20, 1, 0, images[imageTypeWideShotBegin][n], true, opt)
	case resources.SkillWideShotStateMove:
		n := (count / delayWideShot) % len(images[imageTypeWideShotMove])
		next := pos.X + 1
		prev := pos.X - 1
		if direct == config.DirectLeft {
			next, prev = prev, next
		}

		c := count % nextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, pos.X, prev, c, nextStepCount, battlecommon.PanelSize.X)
			dxlib.DrawRotaGraph(view.X+ofsx, view.Y+20, 1, 0, images[imageTypeWideShotMove][n], true, opt)
		}
	}
}
