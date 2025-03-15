package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type tornado struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.Tornado
	drawer  skilldraw.DrawTornado
	animMgr *manager.Manager
}

func newTornado(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *tornado {
	return &tornado{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.Tornado),
		animMgr: animMgr,
	}
}

func (p *tornado) Draw() {
	objPos, targetPos := p.Core.GetPos()
	view := battlecommon.ViewPos(objPos)
	target := battlecommon.ViewPos(targetPos)
	p.drawer.Draw(view, target, p.Core.GetCount(), true)
}

func (p *tornado) Update() (bool, error) {
	return p.Core.Update()
}

func (p *tornado) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *tornado) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
