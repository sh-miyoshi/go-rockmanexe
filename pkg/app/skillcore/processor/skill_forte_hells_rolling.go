package processor

import (
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
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
}

func (p *ForteHellsRolling) Init(skillID int) {
	p.count = 0
	p.currentPos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	p.prevPos = p.currentPos
	p.nextPos = p.currentPos
	if skillID == resources.SkillForteHellsRollingUp {
		p.nextPos.Y--
	} else {
		p.nextPos.Y++
	}
	p.nextPos.X--
	p.curveDirect = 0
}

func (p *ForteHellsRolling) Process() (bool, error) {
	p.count++
	if p.count%forteHellsRollingNextStepCount == 0 {
		p.prevPos = p.currentPos
		p.currentPos = p.nextPos
		if p.currentPos.X < 0 {
			return true, nil
		}

		p.nextPos.X--
		p.nextPos.Y += p.curveDirect

		// 一度だけプレイヤー方向に曲がる
		if p.curveDirect == 0 {
			objType := objanim.ObjTypePlayer
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
