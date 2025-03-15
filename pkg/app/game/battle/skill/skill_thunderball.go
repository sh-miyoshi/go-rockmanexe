package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type thunderBall struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.ThunderBall
	drawer  skilldraw.DrawThunderBall
	animMgr *manager.Manager
}

func newThunderBall(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *thunderBall {
	return &thunderBall{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.ThunderBall),
		animMgr: animMgr,
	}
}

func (p *thunderBall) Draw() {
	prev, current, next := p.Core.GetPos()
	p.drawer.Draw(prev, current, next, p.Core.GetCount(), p.Core.GetNextStepCount())
}

func (p *thunderBall) Update() (bool, error) {
	return p.Core.Update()
}

func (p *thunderBall) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *thunderBall) StopByOwner() {
	// Nothing to do
}
