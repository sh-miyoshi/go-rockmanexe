package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayVulcan = 2 // TODO: 要調整
)

type DrawVulcan struct {
}

func (p *DrawVulcan) Draw(viewPos point.Point, count int) {
	imgNo := 0
	if count > delayVulcan*1 {
		imgNo = (count/(delayVulcan*5))%2 + 1
	}

	// Show body
	dxlib.DrawRotaGraph(viewPos.X+50, viewPos.Y-18, 1, 0, imgVulcan[imgNo], true)
	// Show attack
	if imgNo != 0 {
		if imgNo%2 == 0 {
			dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y-10, 1, 0, imgVulcan[3], true)
		} else {
			dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y-15, 1, 0, imgVulcan[3], true)
		}
	}
}
