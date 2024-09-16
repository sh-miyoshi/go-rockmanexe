package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawForteDarknessOverload struct {
}

func (p *DrawForteDarknessOverload) Draw(pos point.Point, count int) {
	viewPos := battlecommon.ViewPos(pos)
	n := count / 4
	index := n
	if n >= 8 && n < 11 {
		index = n - 3
	} else if n >= 11 && n < 14 {
		index = n - 6
	} else if n >= 14 {
		index = 4 - (n - 14)
	}

	if index >= 0 {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1.0, 0.0, images[imageTypeForteDarknessOverload][index], true)
	}
}
