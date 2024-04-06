package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	heatShotAtkDelay = 15
	heatShotEndCount = 15
)

type HeatShot struct {
	SkillID int
	Arg     skillcore.Argument

	count   int
	targets []point.Point
}

func (p *HeatShot) Process() (bool, error) {
	if p.count == heatShotAtkDelay {
		p.Arg.SoundOn(resources.SEGun)

		pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := point.Point{X: x, Y: pos.Y}
			if objID := p.Arg.GetPanelInfo(target).ObjectID; objID != "" {
				// Hit
				p.Arg.DamageMgr.New(damage.Damage{
					OwnerClientID: p.Arg.OwnerClientID,
					DamageType:    damage.TypeObject,
					Power:         int(p.Arg.Power),
					TargetObjType: p.Arg.TargetType,
					HitEffectType: resources.EffectTypeHeatHit,
					Element:       damage.ElementFire,
					TargetObjID:   objID,
					BigDamage:     true,
				})

				// 誘爆
				p.targets = []point.Point{}
				switch p.SkillID {
				case resources.SkillHeatShot:
					p.targets = append(p.targets, point.Point{X: target.X + 1, Y: target.Y})
				case resources.SkillHeatV:
					p.targets = append(p.targets, point.Point{X: target.X + 1, Y: target.Y - 1})
					p.targets = append(p.targets, point.Point{X: target.X + 1, Y: target.Y + 1})
				case resources.SkillHeatSide:
					p.targets = append(p.targets, point.Point{X: target.X, Y: target.Y - 1})
					p.targets = append(p.targets, point.Point{X: target.X, Y: target.Y + 1})
				}

				for _, t := range p.targets {
					if objID := p.Arg.GetPanelInfo(t).ObjectID; objID != "" {
						p.Arg.DamageMgr.New(damage.Damage{
							OwnerClientID: p.Arg.OwnerClientID,
							DamageType:    damage.TypeObject,
							Power:         int(p.Arg.Power),
							TargetObjType: p.Arg.TargetType,
							HitEffectType: resources.EffectTypeNone,
							Element:       damage.ElementFire,
							TargetObjID:   objID,
							BigDamage:     true,
						})
					}
				}

				break
			}
		}
	}

	p.count++

	if p.count > heatShotEndCount {
		return true, nil
	}
	return false, nil
}

func (p *HeatShot) GetCount() int {
	return p.count
}

func (p *HeatShot) GetDelay() int {
	return heatShotAtkDelay
}

func (p *HeatShot) PopHitTargets() []point.Point {
	if len(p.targets) > 0 {
		res := append([]point.Point{}, p.targets...)
		p.targets = []point.Point{}
		return res
	}
	return []point.Point{}
}
