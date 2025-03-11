package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type boomerang struct {
	ID     string
	Arg    skillcore.Argument
	Core   *processor.Boomerang
	drawer skilldraw.DrawBoomerang
}

func newBoomerang(objID string, arg skillcore.Argument, core skillcore.SkillCore) *boomerang {
	return &boomerang{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.Boomerang),
	}
}

func (p *boomerang) Draw() {
	prev, current, next := p.Core.GetPos()
	p.drawer.Draw(prev, current, next, p.Core.GetCount(), p.Core.GetNextStepCount())
}

func (p *boomerang) Update() (bool, error) {
	return p.Core.Update()
}

func (p *boomerang) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *boomerang) StopByOwner() {
	// Nothing to do after throwing
}
