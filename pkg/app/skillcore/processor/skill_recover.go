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
	Arg skillcore.Argument

	count int
}

func (p *Recover) Update() (bool, error) {
	if p.count == 0 {
		p.Arg.SoundOn(resources.SERecover)
		p.Arg.DamageMgr.New(damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
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
