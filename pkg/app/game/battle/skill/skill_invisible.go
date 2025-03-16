package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type invisible struct {
	ID      string
	Arg     skillcore.Argument
	Core    skillcore.SkillCore
	animMgr *manager.Manager
}

func newInvisible(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *invisible {
	return &invisible{
		ID:      objID,
		Arg:     arg,
		Core:    core,
		animMgr: animMgr,
	}
}

func (p *invisible) Draw() {
}

func (p *invisible) Update() (bool, error) {
	return p.Core.Update()
}

func (p *invisible) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *invisible) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
