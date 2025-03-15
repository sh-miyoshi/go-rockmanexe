package skill

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
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
	animMgr    *manager.Manager
}

func newCountBomb(objID string, arg skillcore.Argument, core skillcore.SkillCore, animMgr *manager.Manager) *countBomb {
	return &countBomb{
		ID:         objID,
		Arg:        arg,
		Core:       core.(*processor.CountBomb),
		objCreated: false,
		animMgr:    animMgr,
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
		if err := obj.Init(p.ID, pm, p.animMgr); err != nil {
			return false, errors.Wrap(err, "count bomb create failed")
		}
		p.animMgr.ObjAnimNew(obj)
		return true, nil
	}

	return false, nil
}

func (p *countBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *countBomb) StopByOwner() {
}
