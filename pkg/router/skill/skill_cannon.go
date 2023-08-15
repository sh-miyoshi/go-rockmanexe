package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
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
		pos := routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID)
		dm := damage.Damage{
			DamageType:    damage.TypeObject,
			OwnerClientID: p.Arg.OwnerClientID,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeCannonHit,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		if p.Arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				if objID := p.Arg.GameInfo.GetPanelInfo(common.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					logger.Debug("Add damage by cannon: %+v", dm)
					routeranim.DamageNew(p.Arg.OwnerClientID, dm)
					break
				}
			}
		} else {
			for x := pos.X - 1; x >= 0; x-- {
				if objID := p.Arg.GameInfo.GetPanelInfo(common.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
					logger.Debug("Add damage by cannon: %+v", dm)
					routeranim.DamageNew(p.Arg.OwnerClientID, dm)
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
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.count,
	}
	switch p.Type {
	case TypeNormalCannon:
		info.AnimType = routeranim.TypeCannonNormal
	case TypeHighCannon:
		info.AnimType = routeranim.TypeCannonHigh
	case TypeMegaCannon:
		info.AnimType = routeranim.TypeCannonMega
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *cannon) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *cannon) GetEndCount() int {
	max := imgCannonBodyNum * delayCannonBody
	if max < imgCannonAtkNum*delayCannonAtk+15 {
		max = imgCannonAtkNum*delayCannonAtk + 15
	}
	return max
}
