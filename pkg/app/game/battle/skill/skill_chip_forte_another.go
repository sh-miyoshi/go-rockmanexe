package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type chipForteAnother struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ChipForteAnother

	drawer skilldraw.DrawForteHellsRolling
}

func newChipForteAnother(objID string, arg skillcore.Argument, core skillcore.SkillCore) *chipForteAnother {
	return &chipForteAnother{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ChipForteAnother),
	}
}

func (p *chipForteAnother) Draw() {
	if p.Core.GetState() == resources.SkillChipForteAnotherStateAttack {
		prev, current, next := p.Core.GetAttackPos()
		// WIP: flip
		p.drawer.Draw(prev, current, next, p.Core.GetAttackCount(), p.Core.GetAttackNextStepCount())
	}
}

func (p *chipForteAnother) Process() (bool, error) {
	end, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	if end {
		field.SetBlackoutCount(0)
		return true, nil
	}
	return false, nil
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
