package skillcore

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	miniBombEndCount = 60
)

type MiniBomb struct {
	arg     Argument
	count   int
	mgrInst *Manager
	target  point.Point
}

func (p *MiniBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		// sound.On(resources.SEBombThrow)// TODO
		pos := p.mgrInst.GetObjectPos(p.arg.OwnerID)
		p.target = point.Point{X: pos.X + 3, Y: pos.Y}
	}

	if p.count == miniBombEndCount {
		pn := p.arg.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		if objID := p.arg.GetPanelInfo(p.target).ObjectID; objID != "" {
			p.mgrInst.damageMgr.New(damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         int(p.arg.Power),
				TargetObjType: p.arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				BigDamage:     true,
				Element:       damage.ElementNone,
				TargetObjID:   objID,
			})
		}
		return true, nil
	}
	return false, nil
}

func (p *MiniBomb) GetCount() int {
	return p.count
}

func (p *MiniBomb) GetEndCount() int {
	return miniBombEndCount
}
