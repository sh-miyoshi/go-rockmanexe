package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawSword struct {
}

func (p *DrawSword) Draw(skillID int, viewPos point.Point, count int, delay int, isReverse bool) {
	opt := dxlib.OptXReverse(isReverse)

	n := (count - 5) / delay
	imgs := getSwordImages(skillID)
	if n >= 0 && n < len(imgs) {
		dxlib.DrawRotaGraph(viewPos.X+math.ReverseIf(100, isReverse), viewPos.Y, 1, 0, imgs[n], true, opt)
	}
}

func getSwordImages(id int) []int {
	n := len(images[imageTypeSword]) / 3

	switch id {
	case resources.SkillSword:
		return images[imageTypeSword][0 : n-1]
	case resources.SkillWideSword, resources.SkillNonEffectWideSword:
		return images[imageTypeSword][2*n:]
	case resources.SkillLongSword:
		return images[imageTypeSword][n : 2*n-1]
	case resources.SkillFighterSword:
		return images[imageTypeFighterSword]
	case resources.SkillDreamSword:
		return images[imageTypeDreamSword]
	}
	return []int{}
}
