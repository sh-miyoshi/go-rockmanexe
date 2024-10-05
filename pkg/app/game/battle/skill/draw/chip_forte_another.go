package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawChipForteAnotherParam struct {
	AttackPrevPos       point.Point
	AttackCurrentPos    point.Point
	AttackNextPos       point.Point
	AttackCount         int
	AttackNextStepCount int
}

type DrawChipForteAnother struct {
	drawer DrawForteHellsRolling
}

func (p *DrawChipForteAnother) Draw(count int, state int, viewPos point.Point, param DrawChipForteAnotherParam) {
	if state == resources.SkillChipForteAnotherStateAppear {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y-20, 1, 0, images[imageTypeForteStand][0], true)
	}

	if state == resources.SkillChipForteAnotherStateAttack {
		n := count / 2
		if n >= len(images[imageTypeForteAtatckHellsRolling]) {
			n = len(images[imageTypeForteAtatckHellsRolling]) - 1
		}

		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y-20, 1, 0, images[imageTypeForteAtatckHellsRolling][n], true)
		p.drawer.Draw(param.AttackPrevPos, param.AttackCurrentPos, param.AttackNextPos, param.AttackCount, param.AttackNextStepCount, true)
	}
}
