package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawAquaman struct {
}

func (p *DrawAquaman) Draw(viewPos point.Point, count int, state int) {
	switch state {
	case resources.SkillAquamanStateInit:
	case resources.SkillAquamanStateAppear:
		const delay = 8
		if count > 20 {
			imgNo := (count / delay) % len(images[imageTypeAquamanCharStand])
			dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y, 1, 0, images[imageTypeAquamanCharStand][imgNo], true, dxlib.OptXReverse(true))
		}
	case resources.SkillAquamanStateCreatePipe:
		imgNo := count
		if imgNo >= len(images[imageTypeAquamanCharCreate]) {
			imgNo = len(images[imageTypeAquamanCharCreate]) - 1
		}
		dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y, 1, 0, images[imageTypeAquamanCharCreate][imgNo], true, dxlib.OptXReverse(true))
	case resources.SkillAquamanStateAttack:
		dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y, 1, 0, images[imageTypeAquamanCharCreate][len(images[imageTypeAquamanCharCreate])-1], true, dxlib.OptXReverse(true))
	}
}
