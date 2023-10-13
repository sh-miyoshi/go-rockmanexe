package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawWaterBomb struct {
	imgBombThrow []int
}

func (p *DrawWaterBomb) Init() {
	p.imgBombThrow = imgBombThrow
}

func (p *DrawWaterBomb) Draw(objPos, targetPos common.Point, count int) {
	imgNo := (count / resources.SKillBombThrowDelay) % len(p.imgBombThrow)
	view := battlecommon.ViewPos(objPos)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := battlecommon.PanelSize.X * (targetPos.X - objPos.X)
	ofsx := size * count / resources.SkillWaterBombEndCount
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if targetPos.Y != objPos.Y {
		size = battlecommon.PanelSize.Y * (targetPos.Y - objPos.Y)
		dy := size * count / resources.SkillWaterBombEndCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, p.imgBombThrow[imgNo], true)

}
