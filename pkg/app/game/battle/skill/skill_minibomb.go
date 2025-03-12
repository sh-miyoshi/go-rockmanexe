package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type miniBomb struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.MiniBomb
	drawer  skilldraw.DrawMiniBomb
	animMgr *manager.Manager
}

func newMiniBomb(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *miniBomb {
	return &miniBomb{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.MiniBomb),
		animMgr: animMgr,
	}
}

func (p *miniBomb) Draw() {
	current, target := p.Core.GetPointParams()
	p.drawer.Draw(current, target, p.Core.GetCount(), p.Core.GetLandCount())
}

func (p *miniBomb) Update() (bool, error) {
	end, err := p.Core.Update()
	if err != nil {
		return false, err
	}

	if eff := p.Core.PopEffect(); eff != nil {
		p.animMgr.EffectAnimNew(effect.Get(eff.Type, eff.Pos, eff.RandRange))
	}
	return end, nil
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *miniBomb) StopByOwner() {
	if p.Core.GetCount() < 5 {
		p.animMgr.AnimDelete(p.ID)
	}
}
