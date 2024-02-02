package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawDreamSword struct {
}

func (p *DrawDreamSword) Draw(viewPos point.Point, count int, delay int) {
	n := (count - 5) / delay
	if n >= 0 && n < len(imgDreamSword) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y, 1, 0, imgDreamSword[n], true)
	}
}
