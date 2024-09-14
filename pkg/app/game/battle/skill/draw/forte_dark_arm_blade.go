package skilldraw

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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
	logger.Debug("Draw ForteDarkArmBlade at %d, %d: %d", viewPos.X, viewPos.Y, n+index)
	dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, images[imageTypeForteDarkArmBlade][n+index], true)
}
