package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayShrimpyAttackBegin = 4
	delayShrimpyAttackMove  = 4
)

type DrawShrimpyAtk struct {
}

// WIP
func (p *DrawShrimpyAtk) Draw(pos point.Point, count int, nextStepCount int, state int) {
	view := battlecommon.ViewPos(pos)

	switch state {
	case resources.SkillShrimpyAttackStateBegin:
		n := (count / delayWideShot)

		if n >= len(imgWideShotBegin) {
			n = len(imgWideShotBegin) - 1
		}
		dxlib.DrawRotaGraph(view.X, view.Y+20, 1, 0, imgWideShotBegin[n], true)
	case resources.SkillShrimpyAttackStateMove:
		n := (count / delayWideShot) % len(imgWideShotMove)
		next := pos.X + 1
		prev := pos.X - 1

		c := count % nextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, pos.X, prev, c, nextStepCount, battlecommon.PanelSize.X)
			dxlib.DrawRotaGraph(view.X+ofsx, view.Y+20, 1, 0, imgWideShotMove[n], true)
		}
	}
}
