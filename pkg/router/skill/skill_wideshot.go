package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type wideShot struct {
	ID  string
	Arg Argument

	state    int
	count    int
	pos      point.Point
	damageID [3]string
}

func newWideShot(arg Argument) *wideShot {
	return &wideShot{
		ID:  arg.AnimObjID,
		Arg: arg,
		pos: routeranim.ObjAnimGetObjPos(arg.OwnerClientID, arg.OwnerObjectID),
	}
}

func (p *wideShot) Draw() {
	// nothing to do at router
}

func (p *wideShot) Process() (bool, error) {
	for _, did := range p.damageID {
		if did != "" {
			if !routeranim.DamageManager(p.Arg.OwnerClientID).Exists(did) && p.count%resources.SkillWideShotPlayerNextStepCount != 0 {
				// attack hit to target
				return true, nil
			}
		}
	}

	switch p.state {
	case resources.SkillWideShotStateBegin:
		if p.count > resources.SkillWideShotEndCount {
			p.state = resources.SkillWideShotStateMove
			p.count = 0
			return false, nil
		}
	case resources.SkillWideShotStateMove:
		if p.count%resources.SkillWideShotPlayerNextStepCount == 0 {
			p.pos.X++

			if p.pos.X >= battlecommon.FieldNum.X || p.pos.X < 0 {
				return true, nil
			}

			for i := -1; i <= 1; i++ {
				y := p.pos.Y + i
				if y < 0 || y >= battlecommon.FieldNum.Y {
					continue
				}

				p.damageID[i+1] = routeranim.DamageNew(p.Arg.OwnerClientID, damage.Damage{
					DamageType:    damage.TypePosition,
					OwnerClientID: p.Arg.OwnerClientID,
					Pos:           point.Point{X: p.pos.X, Y: y},
					Power:         int(p.Arg.Power),
					TTL:           resources.SkillWideShotPlayerNextStepCount,
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

func (p *wideShot) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		AnimType:      routeranim.TypeWideShot,
		ActCount:      p.state*1000 + p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *wideShot) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *wideShot) GetEndCount() int {
	return 0
}
