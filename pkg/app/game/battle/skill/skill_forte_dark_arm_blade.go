package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillForteDarkArmBlade struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteDarkArmBlade
}

func newForteDarkArmBlade(objID string, arg skillcore.Argument, core skillcore.SkillCore) *skillForteDarkArmBlade {
	return &skillForteDarkArmBlade{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ForteDarkArmBlade),
	}
}

func (p *skillForteDarkArmBlade) Draw() {
	// p.drawer.Draw()
}

func (p *skillForteDarkArmBlade) Process() (bool, error) {
	return p.Core.Process()
}

func (p *skillForteDarkArmBlade) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillForteDarkArmBlade) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
