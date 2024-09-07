package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayShrimpyAttackBegin = 4
	delayShrimpyAttackMove  = 4
)

type DrawShrimpyAtk struct {
}

func (p *DrawShrimpyAtk) Draw(pos point.Point, count int, nextStepCount int, state int) {
	view := battlecommon.ViewPos(pos)

	switch state {
	case resources.SkillShrimpyAttackStateBegin:
		n := (count / delayShrimpyAttackBegin)

		if n >= len(images[imageTypeShrimpyAtkBegin]) {
			n = len(images[imageTypeShrimpyAtkBegin]) - 1
		}
		dxlib.DrawRotaGraph(view.X, view.Y+30, 1, 0, images[imageTypeShrimpyAtkBegin][n], true)
	case resources.SkillShrimpyAttackStateMove:
		n := (count / delayShrimpyAttackMove) % len(images[imageTypeShrimpyAtkMove])
		next := pos.X - 1
		prev := pos.X + 1
		c := count % nextStepCount
		ofsx := battlecommon.GetOffset(next, pos.X, prev, c, nextStepCount, battlecommon.PanelSize.X) - battlecommon.PanelSize.X/2
		ofsy := 30 - 20*math.MountainIndex(c, nextStepCount)/nextStepCount
		dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, images[imageTypeShrimpyAtkMove][n], true)
	}
}

func (p *DrawShrimpyAtk) IsDrawEnd(count int, state int) bool {
	switch state {
	case resources.SkillShrimpyAttackStateBegin:
		return count >= len(images[imageTypeShrimpyAtkBegin])*delayShrimpyAttackBegin
	case resources.SkillShrimpyAttackStateMove:
		return count >= len(images[imageTypeShrimpyAtkMove])*delayShrimpyAttackMove
	}
	return false
}
