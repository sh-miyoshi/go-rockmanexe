package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayRecover = 1
)

type DrawRecover struct {
	imgRecover []int
}

func (p *DrawRecover) Init() {
	p.imgRecover = imgRecover
}

func (p *DrawRecover) Draw(viewPos common.Point, count int) {
	n := (count / delayRecover) % len(p.imgRecover)
	if n >= 0 {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, p.imgRecover[n], true)
	}
}
