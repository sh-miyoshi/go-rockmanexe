package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayBubbleShot = 3
)

type DrawBubbleShot struct {
}

func (p *DrawBubbleShot) Draw(viewPos point.Point, count int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)
	n := count / delayBubbleShot

	// Show body
	if n < len(imgBubbleShotBody) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(50, !isPlayer), viewPos.Y-18, 1, 0, imgBubbleShotBody[n], true, opt)
	}

	// Show atk
	n = (count - 4) / delayBubbleShot
	if n >= 0 && n < len(imgBubbleShotAtk) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y-20, 1, 0, imgBubbleShotAtk[n], true, opt)
	}
}
