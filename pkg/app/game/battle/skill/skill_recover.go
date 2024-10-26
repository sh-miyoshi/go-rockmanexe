package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type recover struct {
	ID   string
	Arg  skillcore.Argument
	Core skillcore.SkillCore

	drawer skilldraw.DrawRecover
}

func newRecover(objID string, arg skillcore.Argument, core skillcore.SkillCore) *recover {
	return &recover{
		ID:   objID,
		Arg:  arg,
		Core: core,
	}
}

func (p *recover) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount())
}

func (p *recover) Update() (bool, error) {
	return p.Core.Update()
}

func (p *recover) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *recover) StopByOwner() {
	// Nothing to do
}
