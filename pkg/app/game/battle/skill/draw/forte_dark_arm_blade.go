package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawForteDarkArmBlade struct {
}

func (p *DrawForteDarkArmBlade) Draw(pos point.Point, count int, skillID int) {
	index := 0
	switch skillID {
	case resources.SkillForteDarkArmBladeType2:
		index = 4
	}
	n := count / 4
	if n >= (len(images[imageTypeForteDarkArmBlade]) / 2) {
		n = (len(images[imageTypeForteDarkArmBlade]) / 2) - 1
	}
	viewPos := battlecommon.ViewPos(pos)

	opt := dxlib.DrawRotaGraphOption{}
	if skillID == resources.SkillForteDarkArmBladeType2 {
		t := int32(dxlib.TRUE)
		opt.ReverseXFlag = &t
	}

	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, images[imageTypeForteDarkArmBlade][n+index], true, opt)
}
