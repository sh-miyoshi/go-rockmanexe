package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayHeatShot = 3
)

type DrawHeatShot struct {
}

func (p *DrawHeatShot) Draw(viewPos point.Point, count int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)
	n := count / delayHeatShot

	// Show body
	if n < len(images[imageTypeHeatShotBody]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(50, !isPlayer), viewPos.Y-18, 1, 0, images[imageTypeHeatShotBody][n], true, opt)
	}

	// Show atk
	n = (count - 4) / delayHeatShot
	if n >= 0 && n < len(images[imageTypeHeatShotAtk]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y-20, 1, 0, images[imageTypeHeatShotAtk][n], true, opt)
	}
}
