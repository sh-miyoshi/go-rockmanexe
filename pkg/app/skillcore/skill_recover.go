package skillcore

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
)

const (
	recoverEndCount = 8
)

type Recover struct {
	arg     Argument
	count   int
	mgrInst *Manager
}

func (p *Recover) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		p.mgrInst.damageMgr.New(damage.Damage{
			DamageType:    damage.TypeObject,
			Power:         -int(p.arg.Power),
			TargetObjType: p.arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			Element:       damage.ElementNone,
			TargetObjID:   p.arg.OwnerID,
		})
	}

	if p.count > p.GetEndCount() {
		return true, nil
	}
	return false, nil
}

func (p *Recover) GetCount() int {
	return p.count
}

func (p *Recover) GetEndCount() int {
	return recoverEndCount
}
