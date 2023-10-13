package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawAquamanShot struct {
	img []int
}

func (p *DrawAquamanShot) Init() {
	p.img = imgAquamanShot
}

func (p *DrawAquamanShot) Draw(viewPos, ofs common.Point) {
	dxlib.DrawRotaGraph(viewPos.X+ofs.X, viewPos.Y+ofs.Y, 1, 0, p.img[0], true)
}
