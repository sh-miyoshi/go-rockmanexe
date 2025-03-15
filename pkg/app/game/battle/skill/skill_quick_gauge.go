package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type quickGauge struct {
	ID      string
	Arg     skillcore.Argument
	Core    skillcore.SkillCore
	animMgr *manager.Manager
}

func newQuickGauge(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *quickGauge {
	return &quickGauge{
		ID:      objID,
		Arg:     arg,
		Core:    core,
		animMgr: animMgr,
	}
}

func (p *quickGauge) Draw() {
}

func (p *quickGauge) Update() (bool, error) {
	return p.Core.Update()
}

func (p *quickGauge) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *quickGauge) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
