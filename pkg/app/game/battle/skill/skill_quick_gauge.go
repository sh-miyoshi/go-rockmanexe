package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type quickGauge struct {
	ID   string
	Arg  skillcore.Argument
	Core skillcore.SkillCore
}

func newQuickGauge(objID string, arg skillcore.Argument, core skillcore.SkillCore) *quickGauge {
	return &quickGauge{
		ID:   objID,
		Arg:  arg,
		Core: core,
	}
}

func (p *quickGauge) Draw() {
}

func (p *quickGauge) Process() (bool, error) {
	return p.Core.Process()
}

func (p *quickGauge) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *quickGauge) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
