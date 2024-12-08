package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type spreadGun struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.SpreadGun

	drawer skilldraw.DrawSpreadGun
}

type spreadHit struct {
	ID   string
	Core processor.SpreadHit
}

func newSpreadGun(objID string, arg skillcore.Argument, core skillcore.SkillCore) *spreadGun {
	res := &spreadGun{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.SpreadGun),
	}

	return res
}

func (p *spreadGun) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount(), true)
}

func (p *spreadGun) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopSpreadHits() {
		localanim.AnimNew(&spreadHit{ID: uuid.New().String(), Core: hit})
	}

	return res, nil
}

func (p *spreadGun) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *spreadGun) StopByOwner() {
	if p.Core.GetCount() < 5 {
		localanim.AnimDelete(p.ID)
	}
}

func (p *spreadHit) Draw() {
}

func (p *spreadHit) Update() (bool, error) {
	if p.Core.GetCount() == 1 {
		localanim.AnimNew(effect.Get(resources.EffectTypeSpreadHit, p.Core.Pos, 5))
	}
	return p.Core.Update()
}

func (p *spreadHit) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *spreadHit) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
