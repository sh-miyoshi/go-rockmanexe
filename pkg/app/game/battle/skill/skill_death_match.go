package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type deathMatch struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.DeathMatch
}

func newDeathMatch(objID string, arg skillcore.Argument, core skillcore.SkillCore) *deathMatch {
	return &deathMatch{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.DeathMatch),
	}
}

func (p *deathMatch) Draw() {
}

func (p *deathMatch) Update() (bool, error) {
	end, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	if end {
		field.SetBlackoutCount(0)
		return true, nil
	}
	return false, nil
}

func (p *deathMatch) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *deathMatch) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
