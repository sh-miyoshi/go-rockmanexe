package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayTornadoBody = 4
	delayTornadoAtk  = 2
)

type DrawTornado struct {
}

func (p *DrawTornado) Draw(viewPos, targetPos point.Point, count int) {
	n := (count / delayTornadoBody) % len(imgTornadoBody)
	dxlib.DrawRotaGraph(viewPos.X+48, viewPos.Y-12, 1, 0, imgTornadoBody[n], true)

	n = (count / delayTornadoAtk) % len(imgTornadoAtk)
	dxlib.DrawRotaGraph(targetPos.X, targetPos.Y, 1, 0, imgTornadoAtk[n], true)
}
