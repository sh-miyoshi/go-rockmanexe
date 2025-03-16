package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type crack struct {
	ID      string
	Arg     skillcore.Argument
	Core    skillcore.SkillCore
	animMgr *manager.Manager
}

func newCrack(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *crack {
	return &crack{
		ID:      objID,
		Arg:     arg,
		Core:    core,
		animMgr: animMgr,
	}
}

func (p *crack) Draw() {
}

func (p *crack) Update() (bool, error) {
	return p.Core.Update()
}

func (p *crack) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *crack) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
