package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	wideShotEndCount = 16
)

type WideShotParam struct {
	State         int
	Pos           point.Point
	NextStepCount int
	Direct        int
}

type WideShot struct {
	SkillID int
	Arg     skillcore.Argument

	count    int
	damageID [3]string
	pm       WideShotParam
}

func (p *WideShot) Init(isPlayer bool) {
	if isPlayer {
		p.pm.Direct = config.DirectRight
		p.pm.NextStepCount = 8
	} else {
		p.pm.Direct = config.DirectLeft
		p.pm.NextStepCount = 16
	}
	p.pm.Pos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	p.pm.State = resources.SkillWideShotStateBegin
}

func (p *WideShot) Process() (bool, error) {
	for _, did := range p.damageID {
		if did != "" {
			if !p.Arg.DamageMgr.Exists(did) && p.count%p.pm.NextStepCount != 0 {
				// attack hit to target
				return true, nil
			}
		}
	}

	switch p.pm.State {
	case resources.SkillWideShotStateBegin:
		if p.count == 0 {
			p.Arg.SoundOn(resources.SEWideShot)
		}

		if p.count > wideShotEndCount {
			p.pm.State = resources.SkillWideShotStateMove
			p.count = 0
			return false, nil
		}
	case resources.SkillWideShotStateMove:
		if p.count%p.pm.NextStepCount == 0 {
			if p.pm.Direct == config.DirectRight {
				p.pm.Pos.X++
			} else if p.pm.Direct == config.DirectLeft {
				p.pm.Pos.X--
			}

			if p.pm.Pos.X >= battlecommon.FieldNum.X || p.pm.Pos.X < 0 {
				return true, nil
			}

			for i := -1; i <= 1; i++ {
				y := p.pm.Pos.Y + i
				if y < 0 || y >= battlecommon.FieldNum.Y {
					continue
				}

				p.damageID[i+1] = p.Arg.DamageMgr.New(damage.Damage{
					OwnerClientID: p.Arg.OwnerClientID,
					DamageType:    damage.TypePosition,
					Pos:           point.Point{X: p.pm.Pos.X, Y: y},
					Power:         int(p.Arg.Power),
					TTL:           p.pm.NextStepCount,
					TargetObjType: p.Arg.TargetType,
					HitEffectType: resources.EffectTypeNone,
					BigDamage:     true,
					Element:       damage.ElementWater,
				})
			}
		}
	}

	p.count++
	return false, nil
}

func (p *WideShot) GetCount() int {
	return p.count
}

func (p *WideShot) GetEndCount() int {
	return wideShotEndCount
}

func (p *WideShot) GetParam() WideShotParam {
	return p.pm
}
