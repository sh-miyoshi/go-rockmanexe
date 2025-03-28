package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type AirHockey struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.AirHockey
	animMgr *manager.Manager
	drawer  skilldraw.DrawAirHockey
}

func newAirHockey(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *AirHockey {
	return &AirHockey{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.AirHockey),
		animMgr: animMgr,
	}
}

func (p *AirHockey) Draw() {
	prev, current, next := p.Core.GetPos()
	p.drawer.Draw(prev, current, next, p.Core.GetCount(), p.Core.GetNextStepCount())
}

func (p *AirHockey) Update() (bool, error) {
	return p.Core.Update()
}

func (p *AirHockey) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *AirHockey) StopByOwner() {
	// Nothing to do after throwing
}
