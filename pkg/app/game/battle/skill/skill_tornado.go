package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type tornado struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.Tornado

	drawer skilldraw.DrawTornado
}

func newTornado(objID string, arg skillcore.Argument, core skillcore.SkillCore) *tornado {
	return &tornado{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.Tornado),
	}
}

func (p *tornado) Draw() {
	objPos, targetPos := p.Core.GetPos()
	view := battlecommon.ViewPos(objPos)
	target := battlecommon.ViewPos(targetPos)
	p.drawer.Draw(view, target, p.Core.GetCount())
}

func (p *tornado) Process() (bool, error) {
	return p.Core.Process()
}

func (p *tornado) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *tornado) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
