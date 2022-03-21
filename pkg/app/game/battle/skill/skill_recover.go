package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayRecover = 1
)

type recover struct {
	ID  string
	Arg Argument

	count int
}

func newRecover(objID string, arg Argument) *recover {
	return &recover{
		ID:  objID,
		Arg: arg,
	}
}

func (p *recover) Draw() {
	n := (p.count / delayRecover) % len(imgRecover)
	if n >= 0 {
		pos := objanim.GetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgRecover[n], true)
	}
}

func (p *recover) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SERecover)
		pos := objanim.GetObjPos(p.Arg.OwnerID)
		damage.New(damage.Damage{
			Pos:           pos,
			Power:         -int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeNone,
		})
	}

	p.count++

	if p.count > len(imgRecover)*delayRecover {
		return true, nil
	}
	return false, nil
}

func (p *recover) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeEffect,
	}
}

func (p *recover) AtDelete() {
	if p.Arg.RemoveObject != nil {
		p.Arg.RemoveObject(p.ID)
	}
}

func (p *recover) StopByOwner() {
	// Nothing to do
}
