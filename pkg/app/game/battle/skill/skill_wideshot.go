package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type wideShot struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.WideShot

	drawer skilldraw.DrawWideShot
}

func newWideShot(objID string, arg skillcore.Argument, core skillcore.SkillCore) *wideShot {
	return &wideShot{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.WideShot),
	}
}

func (p *wideShot) Draw() {
	pm := p.Core.GetParam()
	p.drawer.Draw(pm.Pos, p.Core.GetCount(), pm.Direct, p.Arg.TargetType == damage.TargetEnemy, pm.NextStepCount, pm.State)
}

func (p *wideShot) Update() (bool, error) {
	return p.Core.Update()
}

func (p *wideShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *wideShot) StopByOwner() {
	if p.Core.GetParam().State != resources.SkillWideShotStateMove {
		localanim.AnimDelete(p.ID)
	}
}
