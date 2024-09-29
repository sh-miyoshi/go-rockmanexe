package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type cannon struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.Cannon
	SkillID int

	drawer skilldraw.DrawCannon
}

func newCannon(objID string, arg skillcore.Argument, core skillcore.SkillCore, skillID int) *cannon {
	return &cannon{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.Cannon),
		SkillID: skillID,
	}
}

func (p *cannon) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(p.SkillID, view, p.Core.GetCount(), true)
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
