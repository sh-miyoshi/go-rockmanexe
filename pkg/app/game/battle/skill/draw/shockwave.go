package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayPick = 3
)

type DrawShockWave struct {
}

func (p *DrawShockWave) Draw(viewPos point.Point, count int, speed int, direct int) {
	n := (count / speed) % len(images[imageTypeShockWave])
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, images[imageTypeShockWave][n], true, dxlib.OptXReverse(direct == config.DirectLeft))
}

type DrawPick struct {
}

func (p *DrawPick) Draw(viewPos point.Point, count int) {
	n := (count / delayPick)
	if n < len(images[imageTypePick]) {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y-15, 1, 0, images[imageTypePick][n], true)
	}
}
