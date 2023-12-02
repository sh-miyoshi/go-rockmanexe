package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawFlamePillerManager struct {
}

func (p *DrawFlamePillerManager) Draw(viewPos point.Point, count int, skillType int, state int) {
	if skillType == resources.SkillFlamePillarTypeLine && state != resources.SkillFlamePillarStateEnd {
		imageNo := count / 4
		if imageNo >= len(imgFlameLineBody) {
			imageNo = len(imgFlameLineBody) - 1
		}

		// Show body
		dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y-15, 1, 0, imgFlameLineBody[imageNo], true)
	}
}

type DrawFlamePiller struct {
}

func (p *DrawFlamePiller) Draw(viewPos point.Point, count int, state int) {
	n := 0
	switch state {
	case resources.SkillFlamePillarStateWakeup:
		n = count / resources.SkillFlamePillarDelay
		if n >= len(imgFlamePillar) {
			n = len(imgFlamePillar) - 1
		}
	case resources.SkillFlamePillarStateDoing:
		t := (count / resources.SkillFlamePillarDelay) % 2
		n = len(imgFlamePillar) - (t + 1)
	case resources.SkillFlamePillarStateEnd:
		n = len(imgFlamePillar) - (1 + count/resources.SkillFlamePillarDelay)
		if n < 0 {
			n = 0
		}
	}

	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, imgFlamePillar[n], true)

}
