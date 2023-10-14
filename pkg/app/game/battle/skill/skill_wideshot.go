package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type wideShot struct {
	ID            string
	Arg           Argument
	Direct        int
	NextStepCount int

	state    int
	count    int
	pos      common.Point
	damageID [3]string
	drawer   skilldraw.DrawWideShot
}

func newWideShot(objID string, arg Argument) *wideShot {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	direct := common.DirectRight
	nextStep := resources.SkillWideShotPlayerNextStepCount
	if arg.TargetType == damage.TargetPlayer {
		direct = common.DirectLeft
		nextStep = 16
	}

	return &wideShot{
		ID:            objID,
		Arg:           arg,
		Direct:        direct,
		NextStepCount: nextStep,
		pos:           pos,
		state:         resources.SkillWideShotStateBegin,
	}
}

func (p *wideShot) Draw() {
	p.drawer.Draw(p.pos, p.count, p.Direct, p.Arg.TargetType == damage.TargetEnemy, p.NextStepCount, p.state)
}

func (p *wideShot) Process() (bool, error) {
	for _, did := range p.damageID {
		if did != "" {
			if !localanim.DamageManager().Exists(did) && p.count%p.NextStepCount != 0 {
				// attack hit to target
				return true, nil
			}
		}
	}

	switch p.state {
	case resources.SkillWideShotStateBegin:
		if p.count == 0 {
			sound.On(resources.SEWideShot)
		}

		if p.count > resources.SkillWideShotEndCount {
			p.state = resources.SkillWideShotStateMove
			p.count = 0
			return false, nil
		}
	case resources.SkillWideShotStateMove:
		if p.count%p.NextStepCount == 0 {
			if p.Direct == common.DirectRight {
				p.pos.X++
			} else if p.Direct == common.DirectLeft {
				p.pos.X--
			}

			if p.pos.X >= battlecommon.FieldNum.X || p.pos.X < 0 {
				return true, nil
			}

			for i := -1; i <= 1; i++ {
				y := p.pos.Y + i
				if y < 0 || y >= battlecommon.FieldNum.Y {
					continue
				}

				p.damageID[i+1] = localanim.DamageManager().New(damage.Damage{
					DamageType:    damage.TypePosition,
					Pos:           common.Point{X: p.pos.X, Y: y},
					Power:         int(p.Arg.Power),
					TTL:           p.NextStepCount,
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
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *wideShot) StopByOwner() {
	if p.state != resources.SkillWideShotStateMove {
		localanim.AnimDelete(p.ID)
	}
}
