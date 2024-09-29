package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayBubbleShotBody = 2
	delayBubbleShotAtk  = 4
)

type DrawBubbleShot struct {
}

func (p *DrawBubbleShot) Draw(viewPos point.Point, count int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)
	n := count / delayBubbleShotBody

	// Show body
	if n < len(images[imageTypeBubbleShotBody]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(50, !isPlayer), viewPos.Y-18, 1, 0, images[imageTypeBubbleShotBody][n], true, opt)
	}

	// Show atk
	n = (count - 4) / delayBubbleShotAtk
	if n >= 0 && n < len(images[imageTypeBubbleShotAtk]) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y-20, 1, 0, images[imageTypeBubbleShotAtk][n], true, opt)
	}
}
