package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawFlamePillerManager struct {
}

func (p *DrawFlamePillerManager) Draw(viewPos point.Point, count int, showBody bool, pillars []processor.FlamePillerParam, delay int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)

	if showBody {
		imageNo := count / 4
		if imageNo >= len(images[imageTypeFlameLineBody]) {
			imageNo = len(images[imageTypeFlameLineBody]) - 1
		}

		// Show body
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(35, !isPlayer), viewPos.Y-15, 1, 0, images[imageTypeFlameLineBody][imageNo], true, opt)
	}

	for _, pillar := range pillars {
		pos := pillar.Point
		if !isPlayer {
			pos.X = battlecommon.FieldNum.X - pos.X - 1
		}
		view := battlecommon.ViewPos(pos)
		drawPillar(view, pillar.Count, pillar.State, delay)
	}
}

func drawPillar(viewPos point.Point, count int, state int, delay int) {
	n := 0
	switch state {
	case resources.SkillFlamePillarStateWakeup:
		n = count / delay
		if n >= len(images[imageTypeFlamePillar]) {
			n = len(images[imageTypeFlamePillar]) - 1
		}
	case resources.SkillFlamePillarStateDoing:
		t := (count / delay) % 2
		n = len(images[imageTypeFlamePillar]) - (t + 1)
	case resources.SkillFlamePillarStateEnd:
		n = len(images[imageTypeFlamePillar]) - (1 + count/delay)
		if n < 0 {
			n = 0
		}
	}

	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, images[imageTypeFlamePillar][n], true)
}
