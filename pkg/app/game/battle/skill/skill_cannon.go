package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type cannon struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.Cannon
	SkillID int

	drawer  skilldraw.DrawCannon
	animMgr *manager.Manager
}

func newCannon(objID string, arg skillcore.Argument, core skillcore.SkillCore, skillID int, animMgr *manager.Manager) *cannon {
	return &cannon{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.Cannon),
		SkillID: skillID,
		animMgr: animMgr,
	}
}

func (p *cannon) Draw() {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(p.SkillID, view, p.Core.GetCount(), true)
}

func (p *cannon) Update() (bool, error) {
	return p.Core.Update()
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *cannon) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
