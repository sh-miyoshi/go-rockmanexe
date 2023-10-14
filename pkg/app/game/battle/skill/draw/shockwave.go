package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawShockWave struct {
}

func (p *DrawShockWave) Draw(viewPos common.Point, count int, speed int, direct int) {
	n := (count / speed) % len(imgShockWave)
	if direct == common.DirectLeft {
		flag := int32(dxlib.TRUE)
		dxopts := dxlib.DrawRotaGraphOption{
			ReverseXFlag: &flag,
		}
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, imgShockWave[n], true, dxopts)
	} else if direct == common.DirectRight {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, imgShockWave[n], true)
	}
}
