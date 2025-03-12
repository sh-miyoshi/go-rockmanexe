package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type sword struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.Sword
	drawer  skilldraw.DrawSword
	animMgr *manager.Manager
}

func newSword(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *sword {
	return &sword{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.Sword),
		animMgr: animMgr,
	}
}

func (p *sword) Draw() {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(p.Core.GetID(), view, p.Core.GetCount(), p.Core.GetDelay(), p.Arg.IsReverse)
}

func (p *sword) Update() (bool, error) {
	return p.Core.Update()
}

func (p *sword) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *sword) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
