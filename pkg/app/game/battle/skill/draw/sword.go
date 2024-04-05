package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawSword struct {
}

func (p *DrawSword) Draw(swordType int, viewPos point.Point, count int, delay int, isPlayer bool) {
	opt := dxlib.DrawRotaGraphOption{}
	if !isPlayer {
		f := int32(dxlib.TRUE)
		opt.ReverseXFlag = &f
	}

	n := (count - 5) / delay
	if n >= 0 && n < len(imgSword[swordType]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y, 1, 0, imgSword[swordType][n], true, opt)
	}
}
