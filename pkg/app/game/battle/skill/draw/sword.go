package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawSword struct {
}

func (p *DrawSword) Draw(swordType int, viewPos point.Point, count int, delay int) {
	n := (count - 5) / delay
	if n >= 0 && n < len(imgSword[swordType]) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y, 1, 0, imgSword[swordType][n], true)
	}
}
