package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayCannonAtk  = 2
	delayCannonBody = 6
)

type DrawCannon struct {
}

func (p *DrawCannon) Draw(skillID int, viewPos point.Point, count int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)
	length := len(images[imageTypeCannonBody]) / 3
	index := 0
	switch skillID {
	case resources.SkillHighCannon:
		index = length
	case resources.SkillMegaCannon:
		index = length * 2
	}

	n := count / delayCannonBody
	ofs := 48
	if n >= 3 && n < 6 {
		ofs -= 15
	}
	if n >= length {
		n = length - 1
	}

	dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(ofs, !isPlayer), viewPos.Y-12, 1, 0, images[imageTypeCannonBody][n+index], true, opt)

	n = (count - 15) / delayCannonAtk
	if n >= 0 && n < length {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(90, !isPlayer), viewPos.Y-10, 1, 0, images[imageTypeCannonAtk][n+index], true, opt)
	}
}
