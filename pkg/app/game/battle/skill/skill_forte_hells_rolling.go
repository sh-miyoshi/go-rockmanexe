package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillHellsRolling struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteHellsRolling

	drawer skilldraw.DrawForteHellsRolling
}

func newForteHellsRolling(objID string, arg skillcore.Argument, core skillcore.SkillCore) *skillHellsRolling {
	return &skillHellsRolling{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ForteHellsRolling),
	}
}

func (p *skillHellsRolling) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.Core.GetCount())
}

func (p *skillHellsRolling) Process() (bool, error) {
	return false, nil
}

func (p *skillHellsRolling) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillHellsRolling) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
