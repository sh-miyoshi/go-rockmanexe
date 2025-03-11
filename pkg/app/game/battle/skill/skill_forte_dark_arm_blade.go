package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type skillForteDarkArmBlade struct {
	ID      string
	Arg     skillcore.Argument
	Core    *processor.ForteDarkArmBlade
	SkillID int

	drawer skilldraw.DrawForteDarkArmBlade
}

func newForteDarkArmBlade(objID string, arg skillcore.Argument, core skillcore.SkillCore, skillID int) *skillForteDarkArmBlade {
	return &skillForteDarkArmBlade{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.ForteDarkArmBlade),
		SkillID: skillID,
	}
}

func (p *skillForteDarkArmBlade) Draw() {
	p.drawer.Draw(p.Core.GetPos(), p.Core.GetCount(), p.SkillID)
}

func (p *skillForteDarkArmBlade) Update() (bool, error) {
	return p.Core.Update()
}

func (p *skillForteDarkArmBlade) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *skillForteDarkArmBlade) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
