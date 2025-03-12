package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type vulcan struct {
	ID      string
	Arg     skillcore.Argument
	Core    (*processor.Vulcan)
	drawer  skilldraw.DrawVulcan
	animMgr *manager.Manager
}

func newVulcan(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *vulcan {
	return &vulcan{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.Vulcan),
		animMgr: animMgr,
	}
}

func (p *vulcan) Draw() {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(view, p.Core.GetCount(), p.Core.GetDelay(), true)
}

func (p *vulcan) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, eff := range p.Core.PopEffects() {
		p.animMgr.EffectAnimNew(effect.Get(eff.Type, eff.Pos, eff.RandRange))
	}

	return res, nil
}

func (p *vulcan) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *vulcan) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
