package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
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
		routeranim.DamageNew(p.Arg.OwnerClientID, damage.Damage{
			DamageType:    damage.TypeObject,
			OwnerClientID: p.Arg.OwnerClientID,
			Power:         -int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			Element:       damage.ElementNone,
			TargetObjID:   p.Arg.OwnerObjectID,
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
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		DrawType:  anim.DrawTypeEffect,
		ExtraInfo: info.Marshal(),
	}
}

func (p *recover) StopByOwner() {
}

func (p *recover) GetEndCount() int {
	return recoverEndCount
}
