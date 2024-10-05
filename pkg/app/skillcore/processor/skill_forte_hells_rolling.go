package processor

import (
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteHellsRollingNextStepCount = 8
)

type ForteHellsRolling struct {
	Arg skillcore.Argument

	count       int
	prevPos     point.Point
	currentPos  point.Point
	nextPos     point.Point
	curveDirect int
	isPlayer    bool
}

func (p *ForteHellsRolling) Init(skillID int, isPlayer bool) {
	p.count = 0
	p.currentPos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	if isPlayer {
		p.currentPos.X++
	} else {
		p.currentPos.X--
	}
	p.prevPos = p.currentPos
	p.nextPos = p.currentPos
	p.isPlayer = isPlayer
	if skillID == resources.SkillForteHellsRollingUp {
		p.nextPos.Y--
	} else {
		p.nextPos.Y++
	}

	p.curveDirect = 0
}

func (p *ForteHellsRolling) Process() (bool, error) {
	if p.count == 0 {
		p.Arg.DamageMgr.New(damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypePosition,
			Pos:           p.currentPos,
			Power:         int(p.Arg.Power),
			TTL:           forteHellsRollingNextStepCount,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
		})
	}

	p.count++
	if p.count%forteHellsRollingNextStepCount == 0 {
		p.prevPos = p.currentPos
		p.currentPos = p.nextPos
		if p.currentPos.X < 0 || p.currentPos.X >= battlecommon.FieldNum.X {
			return true, nil
		}

		if p.isPlayer {
			p.nextPos.X++
		} else {
			p.nextPos.X--
		}
		p.nextPos.Y += p.curveDirect

		// 一度だけプレイヤー方向に曲がる
		if p.curveDirect == 0 {
			objType := objanim.ObjTypePlayer
			if p.isPlayer {
				objType = objanim.ObjTypeEnemy
			}

			objs := p.Arg.GetObjects(objanim.Filter{ObjType: objType})
			if len(objs) > 0 && math.Abs(objs[0].Pos.X-p.nextPos.X) == 1 {
				p.curveDirect = math.Sign(objs[0].Pos.Y - p.nextPos.Y)
			}
		}

		p.Arg.DamageMgr.New(damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypePosition,
			Pos:           p.currentPos,
			Power:         int(p.Arg.Power),
			TTL:           forteHellsRollingNextStepCount,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
		})
	}
	return false, nil
}

func (p *ForteHellsRolling) GetCount() int {
	return p.count
}

func (p *ForteHellsRolling) GetPos() (prev point.Point, current point.Point, next point.Point) {
	return p.prevPos, p.currentPos, p.nextPos
}

func (p *ForteHellsRolling) GetNextStepCount() int {
	return forteHellsRollingNextStepCount
}
