package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawCountBomb struct{}

func (p *DrawCountBomb) Draw(viewPos common.Point, count int) {
	pm := count * 20
	if pm >= 256 {
		pm = 255
	}
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, pm)
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y+16, 1, 0, imgCountBomb[0], true)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)
}
