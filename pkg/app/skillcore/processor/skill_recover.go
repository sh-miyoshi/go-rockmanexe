package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

const (
	recoverEndCount = 8
)

type Recover struct {
	DamageMgr *damage.DamageManager
	Arg       skillcore.Argument

	count int
}

func (p *Recover) Process() (bool, error) {
	if p.count == 0 {
		p.DamageMgr.New(damage.Damage{
			DamageType:    damage.TypeObject,
			Power:         -int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			Element:       damage.ElementNone,
			TargetObjID:   p.Arg.OwnerID,
		})
	}

	p.count++

	if p.count > recoverEndCount {
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
