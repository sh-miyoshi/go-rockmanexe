package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawCountBomb struct{}

func (p *DrawCountBomb) Draw(viewPos point.Point, count int) {
	pm := count * 20
	if pm >= 256 {
		pm = 255
	}
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, pm)
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y+16, 1, 0, images[imageTypeCountBomb][0], true)
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y+20, 1, 0, images[imageTypeCountBombNumber][3], true)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)
}
