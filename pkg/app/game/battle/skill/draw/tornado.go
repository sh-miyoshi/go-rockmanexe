package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayTornadoBody = 4
	delayTornadoAtk  = 2
)

type DrawTornado struct {
}

func (p *DrawTornado) Draw(viewPos, targetPos point.Point, count int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)

	n := (count / delayTornadoBody) % len(images[imageTypeTornadoBody])
	dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(48, !isPlayer), viewPos.Y-12, 1, 0, images[imageTypeTornadoBody][n], true, opt)

	n = (count / delayTornadoAtk) % len(images[imageTypeTornadoAtk])
	dxlib.DrawRotaGraph(targetPos.X, targetPos.Y, 1, 0, images[imageTypeTornadoAtk][n], true)
}
