package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type barrier struct {
	ID   string
	Arg  skillcore.Argument
	Core skillcore.SkillCore
}

func newBarrier(objID string, arg skillcore.Argument, core skillcore.SkillCore) *barrier {
	return &barrier{
		ID:   objID,
		Arg:  arg,
		Core: core,
	}
}

func (p *barrier) Draw() {
	// p.drawer.Draw()
}

func (p *barrier) Update() (bool, error) {
	return p.Core.Update()
}

func (p *barrier) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *barrier) StopByOwner() {
	// Nothing to do
}
