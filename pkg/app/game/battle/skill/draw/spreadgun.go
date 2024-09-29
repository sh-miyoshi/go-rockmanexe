package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const delaySpreadGun = 3

type DrawSpreadGun struct {
}

type DrawSpreadHit struct {
}

func (p *DrawSpreadGun) Draw(viewPos point.Point, count int, isPlayer bool) {
	n := count / delaySpreadGun

	opt := dxlib.OptXReverse(!isPlayer)

	// Show body
	if n < len(images[imageTypeSpreadGunBody]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(50, !isPlayer), viewPos.Y-18, 1, 0, images[imageTypeSpreadGunBody][n], true, opt)
	}

	// Show atk
	n = (count - 4) / delaySpreadGun
	if n >= 0 && n < len(images[imageTypeSpreadGunAtk]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y-20, 1, 0, images[imageTypeSpreadGunAtk][n], true, opt)
	}
}

func (p *DrawSpreadHit) Draw(viewPos point.Point, count int) {
	if count < len(images[imageTypeSpreadHit]) {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y+15, 1, 0, images[imageTypeSpreadHit][count], true)
	}
}
