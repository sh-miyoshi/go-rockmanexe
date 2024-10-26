package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type invisible struct {
	ID   string
	Arg  skillcore.Argument
	Core skillcore.SkillCore
}

func newInvisible(objID string, arg skillcore.Argument, core skillcore.SkillCore) *invisible {
	return &invisible{
		ID:   objID,
		Arg:  arg,
		Core: core,
	}
}

func (p *invisible) Draw() {
}

func (p *invisible) Update() (bool, error) {
	return p.Core.Update()
}

func (p *invisible) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *invisible) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
