package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteDarknessOverloadEndCount = 64
)

type ForteDarknessOverload struct {
	Arg skillcore.Argument

	count int
}

func (p *ForteDarknessOverload) Update() (bool, error) {
	p.count++
	if p.count == p.GetDelay()*3 {
		for x := 0; x < 2; x++ {
			for y := 0; y < battlecommon.FieldNum.Y; y++ {
				pos := point.Point{X: x, Y: y}
				p.Arg.ChangePanelStatus(pos, battlecommon.PanelStatusCrack, 0)
				if objID := p.Arg.GetPanelInfo(pos).ObjectID; objID != "" {
					p.Arg.DamageMgr.New(damage.Damage{
						OwnerClientID: p.Arg.OwnerClientID,
						TargetObjID:   objID,
						DamageType:    damage.TypeObject,
						Power:         int(p.Arg.Power),
						TargetObjType: p.Arg.TargetType,
						HitEffectType: resources.EffectTypeNone,
						StrengthType:  damage.StrengthHigh,
						Element:       damage.ElementNone,
					})
				}
			}
		}
	}

	return p.count >= forteDarknessOverloadEndCount, nil
}

func (p *ForteDarknessOverload) GetCount() int {
	return p.count
}

func (p *ForteDarknessOverload) GetDelay() int {
	return 3
}
