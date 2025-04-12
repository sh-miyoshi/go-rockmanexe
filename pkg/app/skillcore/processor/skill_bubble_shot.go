package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	bubbleShotAtkDelay = 15
	bubbleShotEndCount = 15
)

type BubbleShot struct {
	SkillID int
	Arg     skillcore.Argument

	count   int
	targets []point.Point
}

func (p *BubbleShot) Update() (bool, error) {
	p.Arg.SoundOn(resources.SEBubbleShot)

	if p.count == bubbleShotAtkDelay {
		pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := point.Point{X: x, Y: pos.Y}
			if objID := p.Arg.GetPanelInfo(target).ObjectID; objID != "" {
				p.Arg.SoundOn(resources.SEWaterLanding)
				// Hit
				p.Arg.DamageMgr.New(damage.Damage{
					OwnerClientID: p.Arg.OwnerClientID,
					DamageType:    damage.TypeObject,
					Power:         int(p.Arg.Power),
					TargetObjType: p.Arg.TargetType,
					HitEffectType: resources.EffectTypeWaterBomb,
					Element:       damage.ElementWater,
					TargetObjID:   objID,
					StrengthType:  damage.StrengthHigh,
				})

				// 誘爆
				p.targets = []point.Point{}
				switch p.SkillID {
				case resources.SkillBubbleShot:
					p.targets = append(p.targets, point.Point{X: target.X + 1, Y: target.Y})
				case resources.SkillBubbleV:
					p.targets = append(p.targets, point.Point{X: target.X + 1, Y: target.Y - 1})
					p.targets = append(p.targets, point.Point{X: target.X + 1, Y: target.Y + 1})
				case resources.SkillBubbleSide:
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
							StrengthType:  damage.StrengthHigh,
						})
					}
				}

				break
			}
		}
	}

	p.count++

	if p.count > bubbleShotEndCount {
		return true, nil
	}
	return false, nil
}

func (p *BubbleShot) GetCount() int {
	return p.count
}

func (p *BubbleShot) GetDelay() int {
	return bubbleShotAtkDelay
}

func (p *BubbleShot) PopHitTargets() []point.Point {
	if len(p.targets) > 0 {
		res := append([]point.Point{}, p.targets...)
		p.targets = []point.Point{}
		return res
	}
	return []point.Point{}
}
