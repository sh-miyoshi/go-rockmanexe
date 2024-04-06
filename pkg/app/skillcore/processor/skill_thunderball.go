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
	thunderBallMaxMoveCount  = 6 // debug
	thunderBallNextStepCount = 80
)

type ThunderBall struct {
	Arg skillcore.Argument

	count     int
	moveCount int
	damageID  string
	pos       point.Point
	next      point.Point
	prev      point.Point
}

func (p *ThunderBall) Init() {
	pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
	x := pos.X + 1
	if p.Arg.TargetType == damage.TargetPlayer {
		x = pos.X - 1
	}

	first := point.Point{X: x, Y: pos.Y}
	p.pos = first
	p.prev = pos
	p.next = first
}

func (p *ThunderBall) Process() (bool, error) {
	if p.count == 0 {
		p.Arg.SoundOn(resources.SEThunderBall)
	}

	if p.count%thunderBallNextStepCount == 2 {
		if p.damageID != "" {
			if !p.Arg.DamageMgr.Exists(p.damageID) {
				// attack hit to target
				return true, nil
			}
		}
	}

	if p.count%thunderBallNextStepCount == 0 {
		t := p.pos
		if p.count != 0 {
			// Update current pos
			p.prev = p.pos
			p.pos = p.next

			p.moveCount++
			if p.moveCount > thunderBallMaxMoveCount {
				return true, nil
			}

			if p.pos.X < 0 || p.pos.X > battlecommon.FieldNum.X || p.pos.Y < 0 || p.pos.Y > battlecommon.FieldNum.Y {
				return true, nil
			}
		}

		pn := p.Arg.GetPanelInfo(p.pos)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		p.damageID = p.Arg.DamageMgr.New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           thunderBallNextStepCount + 1,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
			Element:       damage.ElementElec,
			IsParalyzed:   true,
		})

		// Set next pos
		objType := objanim.ObjTypePlayer
		if p.Arg.TargetType == damage.TargetEnemy {
			objType = objanim.ObjTypeEnemy
		}

		objs := p.Arg.GetObjects(objanim.Filter{ObjType: objType})
		if len(objs) == 0 {
			// no target
			if p.Arg.TargetType == damage.TargetPlayer {
				p.next.X--
			} else {
				p.next.X++
			}
		} else {
			xdif := objs[0].Pos.X - t.X
			ydif := objs[0].Pos.Y - t.Y

			if xdif != 0 || ydif != 0 {
				if math.Abs(xdif) > math.Abs(ydif) {
					// move to x
					p.next.X += (xdif / math.Abs(xdif))
				} else {
					// move to y
					p.next.Y += (ydif / math.Abs(ydif))
				}
			}
		}
	}

	p.count++
	return false, nil
}

func (p *ThunderBall) GetCount() int {
	return p.count
}

func (p *ThunderBall) GetPos() (prev, current, next point.Point) {
	return p.prev, p.pos, p.next
}

func (p *ThunderBall) GetNextStepCount() int {
	return thunderBallNextStepCount
}
