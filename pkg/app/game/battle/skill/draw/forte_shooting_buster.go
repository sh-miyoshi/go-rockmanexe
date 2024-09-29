package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawForteShootingBuster struct {
}

func (p *DrawForteShootingBuster) Draw(pos point.Point, count int, initWait int) {
	viewPos := battlecommon.ViewPos(pos)
	n := count / 4 % len(images[imageTypeForteShootingBuster])
	if count > initWait {
		drawHitArea(pos)
	}
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1.0, 0.0, images[imageTypeForteShootingBuster][n], true)
}
