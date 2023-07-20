package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
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
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, imgRecover[n], true)
	}
}

func (p *recover) Process() (bool, error) {
	if p.count == 0 {
		sound.On(resources.SERecover)
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		localanim.DamageManager().New(damage.Damage{
			Pos:           pos,
			Power:         -int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: battlecommon.EffectTypeNone,
			DamageType:    damage.TypeNone,
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
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *recover) StopByOwner() {
	// Nothing to do
}
