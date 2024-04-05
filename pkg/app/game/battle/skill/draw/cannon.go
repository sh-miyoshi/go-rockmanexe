package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayCannonAtk  = 2
	delayCannonBody = 6
)

type DrawCannon struct {
}

func (p *DrawCannon) Draw(cannonType int, viewPos point.Point, count int, isPlayer bool) {
	opt := dxlib.DrawRotaGraphOption{}
	if !isPlayer {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
	}

	n := count / delayCannonBody
	if n < len(imgCannonBody[cannonType]) {
		ofs := 48
		if n >= 3 {
			ofs -= 15
		}
		if !isPlayer {
			ofs *= -1
		}

		dxlib.DrawRotaGraph(viewPos.X+ofs, viewPos.Y-12, 1, 0, imgCannonBody[cannonType][n], true, opt)
	}

	n = (count - 15) / delayCannonAtk
	if n >= 0 && n < len(imgCannonAtk[cannonType]) {
		ofs := 90
		if !isPlayer {
			ofs *= -1
		}
		dxlib.DrawRotaGraph(viewPos.X+ofs, viewPos.Y-10, 1, 0, imgCannonAtk[cannonType][n], true, opt)
	}
}
