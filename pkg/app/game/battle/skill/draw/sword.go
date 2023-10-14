package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const delaySword = 3

type DrawSword struct {
}

func (p *DrawSword) Draw(swordType int, viewPos common.Point, count int) {
	n := (count - 5) / delaySword
	if n >= 0 && n < len(imgSword[swordType]) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y, 1, 0, imgSword[swordType][n], true)
	}
}
