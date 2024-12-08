package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillAreaSteal struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.AreaSteal

	drawer skilldraw.DrawAreaSteal
}

func newAreaSteal(objID string, arg skillcore.Argument, core skillcore.SkillCore) *skillAreaSteal {
	return &skillAreaSteal{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.AreaSteal),
	}
}

func (p *skillAreaSteal) Draw() {
	p.drawer.Draw(p.Core.GetCount(), p.Core.GetState(), p.Core.GetTargets())
}

func (p *skillAreaSteal) Update() (bool, error) {
	return p.Core.Update()
}

func (p *skillAreaSteal) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillAreaSteal) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
