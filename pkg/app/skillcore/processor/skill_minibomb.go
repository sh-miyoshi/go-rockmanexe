package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	miniBombEndCount = 60
)

type MiniBomb struct {
	GetObjectPos func(objID string) point.Point
	DamageMgr    *damage.DamageManager
	Arg          skillcore.Argument

	count  int
	target point.Point
}

func (p *MiniBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		// sound.On(resources.SEBombThrow)// TODO
		pos := p.GetObjectPos(p.Arg.OwnerID)
		p.target = point.Point{X: pos.X + 3, Y: pos.Y}
	}

	if p.count == miniBombEndCount {
		pn := p.Arg.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		if objID := p.Arg.GetPanelInfo(p.target).ObjectID; objID != "" {
			p.DamageMgr.New(damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
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