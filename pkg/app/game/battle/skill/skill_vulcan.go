package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type vulcan struct {
	ID   string
	Arg  skillcore.Argument
	Core (*processor.Vulcan)

	drawer skilldraw.DrawVulcan
}

func newVulcan(objID string, arg skillcore.Argument, core skillcore.SkillCore) *vulcan {
	return &vulcan{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.Vulcan),
	}
}

func (p *vulcan) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(view, p.Core.GetCount(), p.Core.GetDelay(), true)
}

func (p *vulcan) Process() (bool, error) {
	res, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	for _, eff := range p.Core.PopEffects() {
		localanim.AnimNew(effect.Get(eff.Type, eff.Pos, eff.RandRange))
	}

	return res, nil
}

func (p *vulcan) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *vulcan) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
