package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type heatShot struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.HeatShot

	drawer skilldraw.DrawHeatShot
}

func newHeatShot(objID string, arg skillcore.Argument, core skillcore.SkillCore) *heatShot {
	return &heatShot{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.HeatShot),
	}
}

func (p *heatShot) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount())
}

func (p *heatShot) Process() (bool, error) {
	res, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopHitTargets() {
		localanim.AnimNew(effect.Get(resources.EffectTypeHeatHit, hit, 0))
	}
	return res, nil
}

func (p *heatShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *heatShot) StopByOwner() {
	if p.Core.GetCount() < p.Core.GetDelay() {
		localanim.AnimDelete(p.ID)
	}
}
