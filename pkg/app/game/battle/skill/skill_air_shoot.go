package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type AirShoot struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.AirShoot
	animMgr *manager.Manager
	drawer  skilldraw.DrawAirShoot
}

func newAirShoot(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *AirShoot {
	return &AirShoot{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.AirShoot),
		animMgr: animMgr,
	}
}

func (p *AirShoot) Draw() {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	p.drawer.Draw(pos, p.Core.GetCount())
}

func (p *AirShoot) Update() (bool, error) {
	return p.Core.Update()
}

func (p *AirShoot) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *AirShoot) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
