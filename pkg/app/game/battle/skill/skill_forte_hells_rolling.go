package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillHellsRolling struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteHellsRolling

	drawer skilldraw.DrawForteHellsRolling
}

func newForteHellsRolling(objID string, arg skillcore.Argument, core skillcore.SkillCore) *skillHellsRolling {
	return &skillHellsRolling{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ForteHellsRolling),
	}
}

func (p *skillHellsRolling) Draw() {
	prev, current, next := p.Core.GetPos()
	p.drawer.Draw(prev, current, next, p.Core.GetCount(), p.Core.GetNextStepCount(), false)
}

func (p *skillHellsRolling) Update() (bool, error) {
	return p.Core.Update()
}

func (p *skillHellsRolling) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillHellsRolling) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
