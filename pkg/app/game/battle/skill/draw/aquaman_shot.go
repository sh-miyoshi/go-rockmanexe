package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawAquamanShot struct {
}

func (p *DrawAquamanShot) Draw(viewPos, ofs common.Point) {
	dxlib.DrawRotaGraph(viewPos.X+ofs.X, viewPos.Y+ofs.Y, 1, 0, imgAquamanShot[0], true)
}
