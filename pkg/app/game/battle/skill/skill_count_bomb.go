package skill

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type countBomb struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.CountBomb

	drawer     skilldraw.DrawCountBomb
	objCreated bool
}

func newCountBomb(objID string, arg skillcore.Argument, core skillcore.SkillCore) *countBomb {
	return &countBomb{
		ID:         objID,
		Arg:        arg,
		Core:       core.(*processor.CountBomb),
		objCreated: false,
	}
}

func (p *countBomb) Draw() {
	if !p.objCreated {
		pos := p.Core.GetPos()
		view := battlecommon.ViewPos(pos)
		p.drawer.Draw(view, p.Core.GetCount())
	}
}

func (p *countBomb) Update() (bool, error) {
	end, err := p.Core.Update()
	if err != nil {
		return false, errors.Wrap(err, "failed to process count bomb")
	}

	if end {
		obj := &object.CountBomb{}
		pm := object.ObjectParam{
			Pos:           p.Core.GetPos(),
			HP:            50,
			OnwerCharType: objanim.ObjTypePlayer,
			Power:         int(p.Arg.Power),
		}
		if err := obj.Init(p.ID, pm); err != nil {
			return false, errors.Wrap(err, "count bomb create failed")
		}
		localanim.ObjAnimNew(obj)
		return true, nil
	}

	return false, nil
}

func (p *countBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *countBomb) StopByOwner() {
}
