package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type waterBomb struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.WaterBomb
	drawer  skilldraw.DrawWaterBomb
	animMgr *manager.Manager
}

func newWaterBomb(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *waterBomb {
	return &waterBomb{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.WaterBomb),
		animMgr: animMgr,
	}
}

func (p *waterBomb) Draw() {
	current, target := p.Core.GetPointParams()
	p.drawer.Draw(current, target, p.Core.GetCount(), p.Core.GetLandCount())
}

func (p *waterBomb) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopHits() {
		p.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeWaterBomb, hit, 0))
		field.ChangePanelStatus(hit, battlecommon.PanelStatusCrack, 0)
	}
	return res, nil
}

func (p *waterBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *waterBomb) StopByOwner() {
	if p.Core.GetCount() < 5 {
		p.animMgr.AnimDelete(p.ID)
	}
}
