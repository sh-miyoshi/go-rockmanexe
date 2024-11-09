package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type BambooLance struct {
	Arg skillcore.Argument

	count int
}

func (p *BambooLance) Update() (bool, error) {
	p.count++

	if p.count == 5 {
		dm := damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypePosition,
			Pos:           point.Point{X: battlecommon.FieldNum.X - 1},
			Power:         int(p.Arg.Power),
			TTL:           5,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeBambooHit,
			ShowHitArea:   false,
			StrengthType:  damage.StrengthBack,
			PushLeft:      1,
			Element:       damage.ElementWood,
		}
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			dm.Pos.Y = y
			p.Arg.DamageMgr.New(dm)
		}
	}

	if p.count > 10 {
		return true, nil
	}

	return false, nil
}

func (p *BambooLance) GetCount() int {
	return p.count
}
