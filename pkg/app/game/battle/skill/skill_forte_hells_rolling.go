package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillHellsRolling struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteHellsRolling

	drawer  skilldraw.DrawForteHellsRolling
	animMgr *manager.Manager
}

func newForteHellsRolling(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *skillHellsRolling {
	return &skillHellsRolling{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.ForteHellsRolling),
		animMgr: animMgr,
	}
}

func (p *skillHellsRolling) Draw() {
	prev, current, next := p.Core.GetPos()
	p.drawer.Draw(prev, current, next, p.Core.GetCount(), processor.ForteHellsRollingNextStepCount, false)
}

func (p *skillHellsRolling) Update() (bool, error) {
	return p.Core.Update()
}

func (p *skillHellsRolling) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *skillHellsRolling) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
