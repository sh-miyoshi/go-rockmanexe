package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawAquamanShot struct {
}

func (p *DrawAquamanShot) Draw(viewPos, ofs point.Point) {
	dxlib.DrawRotaGraph(viewPos.X+ofs.X, viewPos.Y+ofs.Y, 1, 0, images[imageTypeAquamanShot][0], true)
}
