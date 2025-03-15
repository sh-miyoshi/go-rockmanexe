package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type heatShot struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.HeatShot
	drawer  skilldraw.DrawHeatShot
	animMgr *manager.Manager
}

func newHeatShot(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *heatShot {
	return &heatShot{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.HeatShot),
		animMgr: animMgr,
	}
}

func (p *heatShot) Draw() {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount(), true)
}

func (p *heatShot) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopHitTargets() {
		p.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeHeatHit, hit, 0))
	}
	return res, nil
}

func (p *heatShot) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *heatShot) StopByOwner() {
	if p.Core.GetCount() < p.Core.GetDelay() {
		p.animMgr.AnimDelete(p.ID)
	}
}
