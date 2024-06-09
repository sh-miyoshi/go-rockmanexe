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

func (p *DrawShrimpyAtk) Draw(pos point.Point, count int, nextStepCount int, state int) {
	view := battlecommon.ViewPos(pos)

	switch state {
	case resources.SkillShrimpyAttackStateBegin:
		n := (count / delayShrimpyAttackBegin)

		if n >= len(imgShrimpyAtkBegin) {
			n = len(imgShrimpyAtkBegin) - 1
		}
		dxlib.DrawRotaGraph(view.X, view.Y+30, 1, 0, imgShrimpyAtkBegin[n], true)
	case resources.SkillShrimpyAttackStateMove:
		n := (count / delayShrimpyAttackMove) % len(imgShrimpyAtkMove)
		next := pos.X - 1
		prev := pos.X + 1
		c := count % nextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, pos.X, prev, c, nextStepCount, battlecommon.PanelSize.X)
			dxlib.DrawRotaGraph(view.X+ofsx, view.Y+30, 1, 0, imgShrimpyAtkMove[n], true)
		}
	}
}

func (p *DrawShrimpyAtk) IsDrawEnd(count int, state int) bool {
	switch state {
	case resources.SkillShrimpyAttackStateBegin:
		return count >= len(imgShrimpyAtkBegin)*delayShrimpyAttackBegin
	case resources.SkillShrimpyAttackStateMove:
		return count >= len(imgShrimpyAtkMove)*delayShrimpyAttackMove
	}
	return false
}
