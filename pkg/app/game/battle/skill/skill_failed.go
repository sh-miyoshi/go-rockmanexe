package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type failed struct {
	ID  string
	Arg skillcore.Argument
}

func newFailed(objID string, arg skillcore.Argument) *failed {
	return &failed{
		ID:  objID,
		Arg: arg,
	}
}

func (p *failed) Draw() {
}

func (p *failed) Process() (bool, error) {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	localanim.AnimNew(effect.Get(resources.EffectTypeFailed, pos, 0))
	sound.On(resources.SEFailed)
	return true, nil
}

func (p *failed) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *failed) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
