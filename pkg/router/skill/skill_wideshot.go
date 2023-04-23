package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const (
	wideShotStateBegin int = iota
	wideShotStateMove
)

const (
	delayWideShot         = 4
	wideShotNextStepCount = 8
)

type wideShot struct {
	ID  string
	Arg Argument

	state    int
	count    int
	pos      common.Point
	damageID [3]string
}

func newWideShot(arg Argument) *wideShot {
	return &wideShot{
		ID:  arg.AnimObjID,
		Arg: arg,
	}
}

func (p *wideShot) Draw() {
	// nothing to do at router
}

func (p *wideShot) Process() (bool, error) {
	for _, did := range p.damageID {
		if did != "" {
			if !damage.Exists(did) && p.count%wideShotNextStepCount != 0 {
				// attack hit to target
				return true, nil
			}
		}
	}

	switch p.state {
	case wideShotStateBegin:
		const (
			imgWideShotBodyNum  = 3
			imgWideShotBeginNum = 4
		)
		max := imgWideShotBodyNum
		if imgWideShotBeginNum > max {
			max = imgWideShotBeginNum
		}
		max *= delayWideShot
		if p.count > max {
			p.state = wideShotStateMove
			p.count = 0
			return false, nil
		}
	case wideShotStateMove:
		if p.count%wideShotNextStepCount == 0 {
			p.pos.X++

			if p.pos.X >= battlecommon.FieldNum.X || p.pos.X < 0 {
				return true, nil
			}

			for i := -1; i <= 1; i++ {
				y := p.pos.Y + i
				if y < 0 || y >= battlecommon.FieldNum.Y {
					continue
				}

				p.damageID[i+1] = damage.New(damage.Damage{
					Pos:           common.Point{X: p.pos.X, Y: y},
					Power:         int(p.Arg.Power),
					TTL:           wideShotNextStepCount,
					TargetType:    p.Arg.TargetType,
					HitEffectType: 0, // TODO: 正しい値をセット
					BigDamage:     true,
					DamageType:    damage.TypeWater,
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
		ActCount:      p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       objanim.GetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *wideShot) StopByOwner() {
	routeranim.Delete(p.Arg.OwnerClientID, p.ID)
}

func (p *wideShot) GetEndCount() int {
	return 0
}
