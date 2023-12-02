package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayHeatShot = 3
)

type DrawHeatShot struct {
}

func (p *DrawHeatShot) Draw(viewPos point.Point, count int) {
	n := count / delayHeatShot

	// Show body
	if n < len(imgHeatShotBody) {
		dxlib.DrawRotaGraph(viewPos.X+50, viewPos.Y-18, 1, 0, imgHeatShotBody[n], true)
	}

	// Show atk
	n = (count - 4) / delayHeatShot
	if n >= 0 && n < len(imgHeatShotAtk) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y-20, 1, 0, imgHeatShotAtk[n], true)
	}
}
