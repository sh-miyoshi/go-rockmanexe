package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type spreadGun struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.SpreadGun
	drawer  skilldraw.DrawSpreadGun
	animMgr *manager.Manager
}

type spreadHit struct {
	ID      string
	Core    processor.SpreadHit
	animMgr *manager.Manager
}

func newSpreadGun(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *spreadGun {
	return &spreadGun{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.SpreadGun),
		animMgr: animMgr,
	}
}

func (p *spreadGun) Draw() {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount(), true)
}

func (p *spreadGun) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopSpreadHits() {
		p.animMgr.EffectAnimNew(&spreadHit{
			ID:      uuid.New().String(),
			Core:    hit,
			animMgr: p.animMgr,
		})
	}

	return res, nil
}

func (p *spreadGun) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *spreadGun) StopByOwner() {
	if p.Core.GetCount() < 5 {
		p.animMgr.AnimDelete(p.ID)
	}
}

func (p *spreadHit) Draw() {
}

func (p *spreadHit) Update() (bool, error) {
	if p.Core.GetCount() == 1 {
		p.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeSpreadHit, p.Core.Pos, 5))
	}
	return p.Core.Update()
}

func (p *spreadHit) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *spreadHit) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
