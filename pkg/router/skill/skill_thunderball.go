package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type thunderBall struct {
	ID           string
	Arg          Argument
	MaxMoveCount int

	count     int
	moveCount int
	damageID  string
	pos       common.Point
	next      common.Point
	prev      common.Point
}

const (
	thunderBallNextStepCount = 80
	delayThunderBall         = 6
)

func newThunderBall(max int, arg Argument) *thunderBall {
	pos := routeranim.ObjAnimGetObjPos(arg.OwnerClientID, arg.OwnerObjectID)
	x := pos.X + 1
	if arg.TargetType == damage.TargetPlayer {
		x = pos.X - 1
	}

	first := common.Point{X: x, Y: pos.Y}
	return &thunderBall{
		ID:           arg.AnimObjID,
		Arg:          arg,
		MaxMoveCount: max,
		pos:          first,
		prev:         pos,
		next:         first,
	}
}

func (p *thunderBall) Draw() {
	// nothing to do at router
}

func (p *thunderBall) Process() (bool, error) {
	if p.count%thunderBallNextStepCount == 2 {
		if p.damageID != "" {
			if !damage.Exists(p.damageID) {
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
			if p.moveCount > p.MaxMoveCount {
				return true, nil
			}

			if p.pos.X < 0 || p.pos.X > battlecommon.FieldNum.X || p.pos.Y < 0 || p.pos.Y > battlecommon.FieldNum.Y {
				return true, nil
			}
		}

		pn := p.Arg.GameInfo.GetPanelInfo(p.pos)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		p.damageID = damage.New(damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           thunderBallNextStepCount + 1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: battlecommon.EffectTypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
			DamageType:    damage.TypeElec,
		})

		// Set next pos
		objType := objanim.ObjTypePlayer
		if p.Arg.TargetType == damage.TargetEnemy {
			objType = objanim.ObjTypeEnemy
		}

		objs := routeranim.ObjAnimGetObjs(p.Arg.OwnerClientID, objanim.Filter{ObjType: objType})
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
				if common.Abs(xdif) > common.Abs(ydif) {
					// move to x
					p.next.X += (xdif / common.Abs(xdif))
				} else {
					// move to y
					p.next.Y += (ydif / common.Abs(ydif))
				}
			}
		}
	}

	p.count++
	return false, nil
}

func (p *thunderBall) GetParam() anim.Param {
	// TODO: set offset

	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		AnimType:      routeranim.TypeThunderBall,
		ActCount:      p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *thunderBall) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *thunderBall) GetEndCount() int {
	return 0
}
