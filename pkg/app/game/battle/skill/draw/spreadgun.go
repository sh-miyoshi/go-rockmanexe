package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const delaySpreadGun = 3

type DrawSpreadGun struct {
}

type DrawSpreadHit struct {
}

func (p *DrawSpreadGun) Draw(viewPos point.Point, count int) {
	n := count / delaySpreadGun

	// Show body
	if n < len(imgSpreadGunBody) {
		dxlib.DrawRotaGraph(viewPos.X+50, viewPos.Y-18, 1, 0, imgSpreadGunBody[n], true)
	}

	// Show atk
	n = (count - 4) / delaySpreadGun
	if n >= 0 && n < len(imgSpreadGunAtk) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y-20, 1, 0, imgSpreadGunAtk[n], true)
	}
}

func (p *DrawSpreadHit) Draw(viewPos point.Point, count int) {
	if count < len(imgSpreadHit) {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y+15, 1, 0, imgSpreadHit[count], true)
	}
}
