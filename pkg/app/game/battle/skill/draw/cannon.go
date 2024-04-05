package skilldraw

import (
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

func (p *DrawCannon) Draw(cannonType int, viewPos point.Point, count int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)

	n := count / delayCannonBody
	if n < len(imgCannonBody[cannonType]) {
		ofs := 48
		if n >= 3 {
			ofs -= 15
		}

		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(ofs, !isPlayer), viewPos.Y-12, 1, 0, imgCannonBody[cannonType][n], true, opt)
	}

	n = (count - 15) / delayCannonAtk
	if n >= 0 && n < len(imgCannonAtk[cannonType]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(90, !isPlayer), viewPos.Y-10, 1, 0, imgCannonAtk[cannonType][n], true, opt)
	}
}
