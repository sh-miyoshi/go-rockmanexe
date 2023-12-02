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
	n := (count / speed) % len(imgShockWave)
	if direct == config.DirectLeft {
		flag := int32(dxlib.TRUE)
		dxopts := dxlib.DrawRotaGraphOption{
			ReverseXFlag: &flag,
		}
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, imgShockWave[n], true, dxopts)
	} else if direct == config.DirectRight {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, imgShockWave[n], true)
	}
}

type DrawPick struct {
}

func (p *DrawPick) Draw(viewPos point.Point, count int) {
	n := (count / delayPick)
	if n < len(imgPick) {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y-15, 1, 0, imgPick[n], true)
	}
}
