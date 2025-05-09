package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	cannonEndCount = 34
)

type Cannon struct {
	SkillID int
	Arg     skillcore.Argument

	count int
}

func (p *Cannon) Update() (bool, error) {
	p.count++

	if p.count == 20 {
		p.Arg.SoundOn(resources.SECannon)
		pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
		dm := damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypeObject,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeCannonHit,
			StrengthType:  damage.StrengthHigh,
			Element:       damage.ElementNone,
		}

		if p.Arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				if objID := p.Arg.GetPanelInfo(point.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					p.Arg.DamageMgr.New(dm)
					break
				}
			}
		} else {
			for x := pos.X - 1; x >= 0; x-- {
				if objID := p.Arg.GetPanelInfo(point.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					p.Arg.DamageMgr.New(dm)
					break
				}
			}
		}
	}

	if p.count > cannonEndCount {
		return true, nil
	}
	return false, nil
}

func (p *Cannon) GetCount() int {
	return p.count
}
