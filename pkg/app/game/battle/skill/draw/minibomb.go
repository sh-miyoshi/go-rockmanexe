package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	delayBombThrow = 4
)

type DrawMiniBomb struct {
}

func (p *DrawMiniBomb) Draw(objPos, targetPos point.Point, count int, endCount int) {
	imgNo := (count / delayBombThrow) % len(imgBombThrow)
	view := battlecommon.ViewPos(objPos)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := battlecommon.PanelSize.X * (targetPos.X - objPos.X)
	ofsx := size * count / endCount
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if targetPos.Y != objPos.Y {
		size = battlecommon.PanelSize.Y * (targetPos.Y - objPos.Y)
		dy := size * count / endCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, imgBombThrow[imgNo], true)
}
