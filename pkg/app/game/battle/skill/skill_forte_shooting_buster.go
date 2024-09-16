package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillForteShootingBuster struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteShootingBuster

	drawer skilldraw.DrawForteShootingBuster
}

func newForteShootingBuster(objID string, arg skillcore.Argument, core skillcore.SkillCore) *skillForteShootingBuster {
	return &skillForteShootingBuster{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ForteShootingBuster),
	}
}

func (p *skillForteShootingBuster) Draw() {
	p.drawer.Draw(p.Core.GetPos(), p.Core.GetCount())
}

func (p *skillForteShootingBuster) Process() (bool, error) {
	return p.Core.Process()
}

func (p *skillForteShootingBuster) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillForteShootingBuster) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
