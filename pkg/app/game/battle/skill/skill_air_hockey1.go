package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type AirHockey1 struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.AirHockey1
	animMgr *manager.Manager
}

func newAirHockey1(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *AirHockey1 {
	return &AirHockey1{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.AirHockey1),
		animMgr: animMgr,
	}
}

func (p *AirHockey1) Draw() {
	// TODO: implement draw method
}

func (p *AirHockey1) Update() (bool, error) {
	return p.Core.Update()
}

func (p *AirHockey1) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *AirHockey1) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
