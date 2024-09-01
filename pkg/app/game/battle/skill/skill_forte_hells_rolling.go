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
	p.drawer.Draw(p.Core.GetPos(), p.Core.GetCount(), p.Core.GetNextStepCount())
}

func (p *skillHellsRolling) Process() (bool, error) {
	return p.Core.Process()
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
