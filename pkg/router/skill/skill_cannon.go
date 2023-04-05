package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	TypeNormalCannon int = iota
	TypeHighCannon
	TypeMegaCannon

	TypeCannonMax
)

const (
	delayCannonAtk   = 2
	delayCannonBody  = 6
	imgCannonBodyNum = 5
	imgCannonAtkNum  = 8
)

type cannon struct {
	ID   string
	Type int
	Arg  Argument

	count int
}

func newCannon(cannonType int, arg Argument) *cannon {
	return &cannon{
		ID:   arg.AnimObjID,
		Type: cannonType,
		Arg:  arg,
	}
}

func (p *cannon) Draw() {
	// nothing to do at router
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 20 {
		// Add damage
		pos := objanim.GetObjPos(p.Arg.OwnerID)
		dm := damage.Damage{
			OwnerClientID: p.Arg.OwnerID,
			Pos:           pos,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: 0, // TODO: 正しい値をセット
			BigDamage:     true,
			DamageType:    damage.TypeNone,
		}

		if p.Arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				dm.Pos.X = x
				if p.Arg.GameInfo.GetPanelInfo(common.Point{X: x, Y: dm.Pos.Y}).ObjectID != "" {
					logger.Debug("Add damage by cannon: %+v", dm)
					dm.Pos.X = battlecommon.FieldNum.X - dm.Pos.X - 1
					damage.New(dm)
					break
				}
			}
		} else {
			for x := pos.X - 1; x >= 0; x-- {
				dm.Pos.X = x
				if p.Arg.GameInfo.GetPanelInfo(common.Point{X: x, Y: dm.Pos.Y}).ObjectID != "" {
					logger.Debug("Add damage by cannon: %+v", dm)
					dm.Pos.X = battlecommon.FieldNum.X - dm.Pos.X - 1
					damage.New(dm)
					break
				}
			}
		}
	}

	if p.count > p.GetEndCount() {
		return true, nil
	}
	return false, nil
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *cannon) StopByOwner() {
	anim.Delete(p.ID)
}

func (p *cannon) GetEndCount() int {
	max := imgCannonBodyNum * delayCannonBody
	if max < imgCannonAtkNum*delayCannonAtk+15 {
		max = imgCannonAtkNum*delayCannonAtk + 15
	}
	return max
}
