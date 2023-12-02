package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayTurnadoBody = 4
	delayTurnadoAtk  = 2
)

type DrawTurnado struct {
}

func (p *DrawTurnado) Draw(viewPos, targetPos point.Point, count int) {
	n := (count / delayTurnadoBody) % len(imgTornadoBody)
	dxlib.DrawRotaGraph(viewPos.X+48, viewPos.Y-12, 1, 0, imgTornadoBody[n], true)

	n = (count / delayTurnadoAtk) % len(imgTornadoAtk)
	dxlib.DrawRotaGraph(targetPos.X, targetPos.Y, 1, 0, imgTornadoAtk[n], true)
}
