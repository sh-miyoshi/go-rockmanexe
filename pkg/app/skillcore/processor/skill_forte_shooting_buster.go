package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteShootingBusterEndCount = 16
)

type ForteShootingBuster struct {
	Arg skillcore.Argument

	count int
	pos   point.Point
}

func (p *ForteShootingBuster) Init() {
	p.count = 0
	// WIP: pos
}

func (p *ForteShootingBuster) Process() (bool, error) {
	p.count++
	if p.count == 1 {
		if objID := p.Arg.GetPanelInfo(p.pos).ObjectID; objID != "" {
			p.Arg.DamageMgr.New(damage.Damage{
				OwnerClientID: p.Arg.OwnerClientID,
				TargetObjID:   objID,
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				BigDamage:     true,
				Element:       damage.ElementNone,
			})
		}
	}
	return p.count >= forteShootingBusterEndCount, nil
}

func (p *ForteShootingBuster) GetCount() int {
	return p.count
}

func (p *ForteShootingBuster) GetPos() point.Point {
	return p.pos
}
