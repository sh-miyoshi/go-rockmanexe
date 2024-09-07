package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawDreamSword struct {
}

func (p *DrawDreamSword) Draw(viewPos point.Point, count int, delay int) {
	n := (count - 5) / delay
	if n >= 0 && n < len(images[imageTypeDreamSword]) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y, 1, 0, images[imageTypeDreamSword][n], true)
	}
}
