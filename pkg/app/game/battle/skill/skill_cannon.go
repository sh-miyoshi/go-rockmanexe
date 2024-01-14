package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type cannon struct {
	ID   string
	Type int
	Arg  skillcore.Argument
	Core skillcore.SkillCore

	drawer skilldraw.DrawCannon
}

func newCannon(objID string, cannonType int, arg skillcore.Argument, core skillcore.SkillCore) *cannon {
	return &cannon{
		ID:   objID,
		Type: cannonType,
		Arg:  arg,
		Core: core,
	}
}

func (p *cannon) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(p.Type, view, p.Core.GetCount())
}

func (p *cannon) Process() (bool, error) {
	return p.Core.Process()
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *cannon) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
