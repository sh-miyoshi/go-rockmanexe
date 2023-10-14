package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawDreamSword struct {
}

func (p *DrawDreamSword) Draw(viewPos common.Point, count int) {
	n := (count - 5) / resources.SkillSwordDelay
	if n >= 0 && n < len(imgDreamSword) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y, 1, 0, imgDreamSword[n], true)
	}
}
