package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawSword struct {
}

func (p *DrawSword) Draw(skillID int, viewPos point.Point, count int, delay int, isPlayer bool) {
	opt := dxlib.OptXReverse(!isPlayer)

	n := (count - 5) / delay
	imgs := getSwordImages(skillID)
	if n >= 0 && n < len(imgs) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, !isPlayer), viewPos.Y, 1, 0, imgs[n], true, opt)
	}
}

func getSwordImages(id int) []int {
	switch id {
	case resources.SkillSword:
		return imgSword[0]
	case resources.SkillWideSword:
		return imgSword[1]
	case resources.SkillLongSword:
		return imgSword[2]
	case resources.SkillDreamSword:
		return imgDreamSword
	}
	return []int{}
}
