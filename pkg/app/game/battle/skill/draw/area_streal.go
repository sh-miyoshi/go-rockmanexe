package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayStealHit = 2
)

type DrawAreaSteal struct {
}

func (p *DrawAreaSteal) Draw(count int, state int, targets []point.Point) {
	switch state {
	case resources.SkillAreaStealStateBlackout:
	case resources.SkillAreaStealStateActing:
		ofs := count*4 - 30
		ino := count / 3
		if ino >= len(images[imageTypeAreaStealMain]) {
			ino = len(images[imageTypeAreaStealMain]) - 1
		}

		for _, target := range targets {
			view := battlecommon.ViewPos(target)
			dxlib.DrawRotaGraph(view.X, view.Y+ofs, 1, 0, images[imageTypeAreaStealMain][ino], true)
		}
	case resources.SkillAreaStealStateHit:
		ino := count / delayStealHit
		if ino >= len(images[imageTypeAreaStealPanel]) {
			ino = len(images[imageTypeAreaStealPanel]) - 1
		}

		for _, target := range targets {
			view := battlecommon.ViewPos(target)
			dxlib.DrawRotaGraph(view.X, view.Y+30, 1, 0, images[imageTypeAreaStealPanel][ino], true)
		}
	}
}
