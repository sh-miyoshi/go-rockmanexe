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
	ID      string
	Arg     skillcore.Argument
	Core    *processor.Sword
	SkillID int

	drawer skilldraw.DrawSword
}

func newSword(objID string, skillID int, arg skillcore.Argument, core skillcore.SkillCore) *sword {
	return &sword{
		ID:      objID,
		Arg:     arg,
		SkillID: skillID,
		Core:    core.(*processor.Sword),
	}
}

func (p *sword) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(p.Core.GetSwordType(), view, p.Core.GetCount())
}

func (p *sword) Process() (bool, error) {
	return p.Core.Process()
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
