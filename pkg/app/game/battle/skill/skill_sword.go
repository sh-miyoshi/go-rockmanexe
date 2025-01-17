package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type sword struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.Sword

	drawer skilldraw.DrawSword
}

func newSword(objID string, arg skillcore.Argument, core skillcore.SkillCore) *sword {
	return &sword{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.Sword),
	}
}

func (p *sword) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(p.Core.GetID(), view, p.Core.GetCount(), p.Core.GetDelay(), p.Arg.IsReverse)
}

func (p *sword) Update() (bool, error) {
	return p.Core.Update()
}

func (p *sword) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *sword) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
