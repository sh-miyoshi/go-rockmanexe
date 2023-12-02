package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayRecover = 1
)

type DrawRecover struct {
}

func (p *DrawRecover) Draw(viewPos point.Point, count int) {
	n := (count / delayRecover) % len(imgRecover)
	if n >= 0 {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, imgRecover[n], true)
	}
}
