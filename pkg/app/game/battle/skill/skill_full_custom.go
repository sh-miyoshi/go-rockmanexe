package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type fullCustom struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.FullCustom
}

func newFullCustom(objID string, arg skillcore.Argument, core skillcore.SkillCore) *fullCustom {
	return &fullCustom{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.FullCustom),
	}
}

func (p *fullCustom) Draw() {
}

func (p *fullCustom) Process() (bool, error) {
	return p.Core.Process()
}

func (p *fullCustom) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *fullCustom) StopByOwner() {
}
