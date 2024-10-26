package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
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
}

func (p *panelReturn) Process() (bool, error) {
	end, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	if end {
		field.SetBlackoutCount(0)
		return true, nil
	}
	return false, nil
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
