package skillcore

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	cannonEndCount = 31 // imgAtkNum*delayAtk + 15
)

type Cannon struct {
	arg     Argument
	count   int
	mgrInst *Manager
}

func (p *Cannon) Process() (bool, error) {
	p.count++

	if p.count == 20 {
		pos := p.mgrInst.GetObjectPos(p.arg.OwnerID)
		dm := damage.Damage{
			DamageType:    damage.TypeObject,
			Power:         int(p.arg.Power),
			TargetObjType: p.arg.TargetType,
			HitEffectType: resources.EffectTypeCannonHit,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		if p.arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				if objID := p.arg.GetPanelInfo(point.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					p.mgrInst.damageMgr.New(dm)
					break
				}
			}
		} else {
			for x := pos.X - 1; x >= 0; x-- {
				if objID := p.arg.GetPanelInfo(point.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					p.mgrInst.damageMgr.New(dm)
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

func (p *Cannon) GetEndCount() int {
	return cannonEndCount
}
