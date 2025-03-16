package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type panelReturn struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.PanelReturn
	animMgr *manager.Manager
}

func newPanelReturn(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *panelReturn {
	return &panelReturn{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.PanelReturn),
		animMgr: animMgr,
	}
}

func (p *panelReturn) Draw() {
}

func (p *panelReturn) Update() (bool, error) {
	end, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	if end {
		field.SetBlackoutCount(0)
		return true, nil
	}
	return false, nil
}

func (p *panelReturn) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *panelReturn) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
