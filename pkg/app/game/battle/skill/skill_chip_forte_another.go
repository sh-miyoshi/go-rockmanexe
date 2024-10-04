package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type chipForteAnother struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ChipForteAnother
}

func newChipForteAnother(objID string, arg skillcore.Argument, core skillcore.SkillCore) *chipForteAnother {
	return &chipForteAnother{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ChipForteAnother),
	}
}

func (p *chipForteAnother) Draw() {
	// p.drawer.Draw()
}

func (p *chipForteAnother) Process() (bool, error) {
	return p.Core.Process()
}

func (p *chipForteAnother) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *chipForteAnother) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
