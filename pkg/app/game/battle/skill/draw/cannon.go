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

func (p *DrawCannon) Draw(cannonType int, viewPos point.Point, count int) {
	n := count / delayCannonBody
	if n < len(imgCannonBody[cannonType]) {
		if n >= 3 {
			viewPos.X -= 15
		}

		dxlib.DrawRotaGraph(viewPos.X+48, viewPos.Y-12, 1, 0, imgCannonBody[cannonType][n], true)
	}

	n = (count - 15) / delayCannonAtk
	if n >= 0 && n < len(imgCannonAtk[cannonType]) {
		dxlib.DrawRotaGraph(viewPos.X+90, viewPos.Y-10, 1, 0, imgCannonAtk[cannonType][n], true)
	}
}
