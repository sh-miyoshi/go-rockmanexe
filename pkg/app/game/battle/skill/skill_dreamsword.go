package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type dreamSword struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.Sword

	drawer skilldraw.DrawDreamSword
}

func newDreamSword(objID string, arg skillcore.Argument, core skillcore.SkillCore) *dreamSword {
	return &dreamSword{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.Sword),
	}
}

func (p *dreamSword) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(view, p.Core.GetCount(), p.Core.GetDelay())
}

func (p *dreamSword) Process() (bool, error) {
	return p.Core.Process()
}

func (p *dreamSword) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *dreamSword) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
