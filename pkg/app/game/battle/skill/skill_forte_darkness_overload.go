package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type skillForteDarknessOverload struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteDarknessOverload

	drawer skilldraw.DrawForteDarknessOverload
}

func newForteDarknessOverload(objID string, arg skillcore.Argument, core skillcore.SkillCore) *skillForteDarknessOverload {
	return &skillForteDarknessOverload{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.ForteDarknessOverload),
	}
}

func (p *skillForteDarknessOverload) Draw() {
	pos := point.Point{X: 0, Y: 1}
	p.drawer.Draw(pos, p.Core.GetCount())
}

func (p *skillForteDarknessOverload) Process() (bool, error) {
	return p.Core.Process()
}

func (p *skillForteDarknessOverload) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillForteDarknessOverload) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
