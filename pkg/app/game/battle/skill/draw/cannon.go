package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayCannonAtk  = 2
	delayCannonBody = 6
)

type DrawCannon struct {
	imgBody [resources.SkillTypeCannonMax][]int
	imgAtk  [resources.SkillTypeCannonMax][]int
}

func (p *DrawCannon) Init() {
	p.imgAtk = imgCannonAtk
	p.imgBody = imgCannonBody
}

func (p *DrawCannon) Draw(cannonType int, viewPos common.Point, count int) {
	n := count / delayCannonBody
	if n < len(p.imgBody[cannonType]) {
		if n >= 3 {
			viewPos.X -= 15
		}

		dxlib.DrawRotaGraph(viewPos.X+48, viewPos.Y-12, 1, 0, p.imgBody[cannonType][n], true)
	}

	n = (count - 15) / delayCannonAtk
	if n >= 0 && n < len(p.imgAtk[cannonType]) {
		dxlib.DrawRotaGraph(viewPos.X+90, viewPos.Y-10, 1, 0, p.imgAtk[cannonType][n], true)
	}
}
