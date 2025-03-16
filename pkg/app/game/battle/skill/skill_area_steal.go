package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillAreaSteal struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.AreaSteal

	drawer  skilldraw.DrawAreaSteal
	animMgr *manager.Manager
}

func newAreaSteal(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *skillAreaSteal {
	return &skillAreaSteal{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.AreaSteal),
		animMgr: animMgr,
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
		ObjID: p.ID,
	}
}

func (p *skillAreaSteal) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
