package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type skillForteDarknessOverload struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.ForteDarknessOverload

	drawer  skilldraw.DrawForteDarknessOverload
	animMgr *manager.Manager
}

func newForteDarknessOverload(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *skillForteDarknessOverload {
	return &skillForteDarknessOverload{
		ID:      objID,
		Arg:     arg,
		Core:    core.(*processor.ForteDarknessOverload),
		animMgr: animMgr,
	}
}

func (p *skillForteDarknessOverload) Draw() {
	pos := point.Point{X: 0, Y: 1}
	p.drawer.Draw(pos, p.Core.GetCount(), p.Core.GetDelay())
}

func (p *skillForteDarknessOverload) Update() (bool, error) {
	return p.Core.Update()
}

func (p *skillForteDarknessOverload) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *skillForteDarknessOverload) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
