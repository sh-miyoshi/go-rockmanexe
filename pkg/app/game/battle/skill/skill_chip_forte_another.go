package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

type chipForteAnother struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ChipForteAnother

	drawer skilldraw.DrawChipForteAnother
}

func newChipForteAnother(objID string, arg skillcore.Argument, core skillcore.SkillCore) *chipForteAnother {
	return &chipForteAnother{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ChipForteAnother),
	}
}

func (p *chipForteAnother) Draw() {
	state := p.Core.GetState()
	pm := skilldraw.DrawChipForteAnotherParam{}
	if state == resources.SkillChipForteAnotherStateAttack {
		prev, current, next := p.Core.GetAttackPos()
		pm.AttackPrevPos = prev
		pm.AttackCurrentPos = current
		pm.AttackNextPos = next
		pm.AttackCount = p.Core.GetAttackCount()
		pm.AttackNextStepCount = p.Core.GetAttackNextStepCount()
	}
	view := battlecommon.ViewPos(p.Core.GetPos())
	p.drawer.Draw(p.Core.GetCount(), state, view, pm)
}

func (p *chipForteAnother) Process() (bool, error) {
	end, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	if end {
		field.SetBlackoutCount(0)
		logger.Debug("Forte Another Skill End")
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
