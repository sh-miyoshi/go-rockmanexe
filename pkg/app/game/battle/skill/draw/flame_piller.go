package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawFlamePillerManager struct {
}

func (p *DrawFlamePillerManager) Draw(viewPos point.Point, count int, showBody bool, pillars []processor.FlamePillerParam, delay int) {
	if showBody {
		imageNo := count / 4
		if imageNo >= len(imgFlameLineBody) {
			imageNo = len(imgFlameLineBody) - 1
		}

		// Show body
		dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y-15, 1, 0, imgFlameLineBody[imageNo], true)
	}

	for _, pillar := range pillars {
		view := battlecommon.ViewPos(pillar.Point)
		drawPillar(view, pillar.Count, pillar.State, delay)
	}
}

func drawPillar(viewPos point.Point, count int, state int, delay int) {
	n := 0
	switch state {
	case resources.SkillFlamePillarStateWakeup:
		n = count / delay
		if n >= len(imgFlamePillar) {
			n = len(imgFlamePillar) - 1
		}
	case resources.SkillFlamePillarStateDoing:
		t := (count / delay) % 2
		n = len(imgFlamePillar) - (t + 1)
	case resources.SkillFlamePillarStateEnd:
		n = len(imgFlamePillar) - (1 + count/delay)
		if n < 0 {
			n = 0
		}
	}

	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, imgFlamePillar[n], true)
}
