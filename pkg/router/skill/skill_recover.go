package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/effect"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const (
	recoverEndCount = 8
)

type recover struct {
	ID    string
	Arg   Argument
	count int
}

func newRecover(arg Argument) *recover {
	return &recover{
		ID:  arg.AnimObjID,
		Arg: arg,
	}
}

func (p *recover) Draw() {
	// nothing to do at router
}

func (p *recover) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		pos := objanim.GetObjPos(p.Arg.OwnerObjectID)
		damage.New(damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			Pos:           pos,
			Power:         -int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeNone,
			DamageType:    damage.TypeNone,
		})
	}

	if p.count > p.GetEndCount() {
		return true, nil
	}
	return false, nil
}

func (p *recover) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeRecover,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       objanim.GetObjPos(p.Arg.OwnerObjectID),
		DrawType:  anim.DrawTypeEffect,
		ExtraInfo: info.Marshal(),
	}
}

func (p *recover) StopByOwner() {
}

func (p *recover) GetEndCount() int {
	return recoverEndCount
}
