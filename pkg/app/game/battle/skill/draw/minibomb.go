package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayMiniBombThrow = 4
)

type DrawMiniBomb struct {
	imgBombThrow []int
}

func (p *DrawMiniBomb) Init() error {
	path := common.ImagePath + "battle/skill/"

	fname := path + "ミニボム.png"
	p.imgBombThrow = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 40, 30, p.imgBombThrow); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	return nil
}

func (p *DrawMiniBomb) End() {
	for i := 0; i < len(p.imgBombThrow); i++ {
		dxlib.DeleteGraph(p.imgBombThrow[i])
	}
	p.imgBombThrow = []int{}
}

func (p *DrawMiniBomb) Draw(objPos, targetPos common.Point, count int) {
	imgNo := (count / delayMiniBombThrow) % len(p.imgBombThrow)
	view := battlecommon.ViewPos(objPos)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	// y = (4 * ymax / d^2)x^2 + (4 * ymax / d)x
	size := battlecommon.PanelSize.X * (targetPos.X - objPos.X)
	ofsx := size * count / resources.SkillMiniBombEndCount
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if targetPos.Y != objPos.Y {
		size = battlecommon.PanelSize.Y * (targetPos.Y - objPos.Y)
		dy := size * count / resources.SkillMiniBombEndCount
		ofsy += dy
	}

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+ofsy, 1, 0, p.imgBombThrow[imgNo], true)
}
