package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawVulcan struct {
}

func (p *DrawVulcan) Draw(viewPos point.Point, count int, delay int, isPlayer bool) {
	imgNo := 0
	if count > delay*1 {
		imgNo = (count/(delay*5))%2 + 1
	}
	opt := dxlib.OptXReverse(!isPlayer)

	// Show body
	dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(50, !isPlayer), viewPos.Y-18, 1, 0, images[imageTypeVulcan][imgNo], true, opt)
	// Show attack
	if imgNo != 0 {
		if imgNo%2 == 0 {
			dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y-10, 1, 0, images[imageTypeVulcan][3], true, opt)
		} else {
			dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y-15, 1, 0, images[imageTypeVulcan][3], true, opt)
		}
	}
}
