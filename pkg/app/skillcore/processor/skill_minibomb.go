package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	miniBombLandCount = 60
)

type MiniBomb struct {
	Arg skillcore.Argument

	count  int
	pos    point.Point
	target point.Point
	effect *resources.EffectParam
}

func (p *MiniBomb) Init() {
	p.pos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	p.target = point.Point{X: p.pos.X + 3, Y: p.pos.Y}
}

func (p *MiniBomb) Update() (bool, error) {
	p.count++

	if p.count == 1 {
		p.Arg.SoundOn(resources.SEBombThrow)
	}

	if p.count == miniBombLandCount {
		pn := p.Arg.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		if objID := p.Arg.GetPanelInfo(p.target).ObjectID; objID != "" {
			p.Arg.DamageMgr.New(damage.Damage{
				OwnerClientID: p.Arg.OwnerClientID,
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				StrengthType:  damage.StrengthHigh,
				Element:       damage.ElementNone,
				TargetObjID:   objID,
			})
		}
		p.effect = &resources.EffectParam{
			Type:      resources.EffectTypeExplode,
			Pos:       p.target,
			RandRange: 0,
		}

		return true, nil
	}
	return false, nil
}

func (p *MiniBomb) GetCount() int {
	return p.count
}

func (p *MiniBomb) GetLandCount() int {
	return miniBombLandCount
}

func (p *MiniBomb) GetPointParams() (current, target point.Point) {
	return p.pos, p.target
}

func (p *MiniBomb) PopEffect() *resources.EffectParam {
	if p.effect != nil {
		res := &resources.EffectParam{}
		*res = *p.effect
		p.effect = nil
		return res
	}
	return nil
}
