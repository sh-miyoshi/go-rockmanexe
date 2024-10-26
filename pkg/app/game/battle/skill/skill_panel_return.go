package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type panelReturn struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.PanelReturn
}

func newPanelReturn(objID string, arg skillcore.Argument, core skillcore.SkillCore) *panelReturn {
	return &panelReturn{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.PanelReturn),
	}
}

func (p *panelReturn) Draw() {
	// p.drawer.Draw()
}

func (p *panelReturn) Process() (bool, error) {
	return p.Core.Process()
}

func (p *panelReturn) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *panelReturn) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
