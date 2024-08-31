package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawForteHellsRolling struct {
}

func (p *DrawForteHellsRolling) Draw(viewPos point.Point, count int) {
	n := count / 10 % len(imgForteHellsRolling)
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1.0, 0.0, imgForteHellsRolling[n], true)
}
