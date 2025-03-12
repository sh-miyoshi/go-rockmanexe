package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type failed struct {
	ID      string
	Arg     skillcore.Argument
	animMgr *manager.Manager
}

func newFailed(objID string, arg skillcore.Argument, animMgr *manager.Manager) *failed {
	return &failed{
		ID:      objID,
		Arg:     arg,
		animMgr: animMgr,
	}
}

func (p *failed) Draw() {
}

func (p *failed) Update() (bool, error) {
	pos := p.animMgr.ObjAnimGetObjPos(p.Arg.OwnerID)
	p.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeFailed, pos, 0))
	sound.On(resources.SEFailed)
	return true, nil
}

func (p *failed) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *failed) StopByOwner() {
	p.animMgr.AnimDelete(p.ID)
}
