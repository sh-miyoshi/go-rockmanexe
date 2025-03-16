package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type FullCustom struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.FullCustom
	animMgr *manager.Manager
}

func newFullCustom(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *FullCustom {
	return &FullCustom{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.FullCustom),
		animMgr: animMgr,
	}
}

func (p *FullCustom) Draw() {
	// TODO: implement draw method
}

func (p *FullCustom) Update() (bool, error) {
	return p.Core.Update()
}

func (p *FullCustom) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *FullCustom) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
