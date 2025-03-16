package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type flamePillarManager struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.FlamePillarManager
	drawer  skilldraw.DrawFlamePillerManager
	animMgr *manager.Manager
}

func newFlamePillar(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *flamePillarManager {
	return &flamePillarManager{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.FlamePillarManager),
		animMgr: animMgr,
	}
}

func (p *flamePillarManager) Draw() {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount(), p.Core.IsShowBody(), p.Core.GetPillars(), p.Core.GetDelay(), true)
}

func (p *flamePillarManager) Update() (bool, error) {
	return p.Core.Update()
}

func (p *flamePillarManager) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *flamePillarManager) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
