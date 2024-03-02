package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type crack struct {
	ID   string
	Arg  skillcore.Argument
	Core skillcore.SkillCore
}

func newCrack(objID string, arg skillcore.Argument, core skillcore.SkillCore) *crack {
	return &crack{
		ID:   objID,
		Arg:  arg,
		Core: core,
	}
}

func (p *crack) Draw() {
}

func (p *crack) Process() (bool, error) {
	return p.Core.Process()
}

func (p *crack) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *crack) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
