package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	boomerangActTypeClockwise = iota
	boomerangActTypeCounterClockwise
	boomerangActTypeLeftStraight
	boomerangActTypeRightStraight
)

const (
	boomerangNextStepCount = 6
)

type Boomerang struct {
	Arg skillcore.Argument

	actType int
	count   int
	turnNum int
	pos     point.Point
	next    point.Point
	prev    point.Point
}

func (p *Boomerang) Init() {
	// TODO: actTypeは別で渡すようにする
	// TODO(直線移動)
	sx := 0
	sy := battlecommon.FieldNum.Y - 1
	act := boomerangActTypeCounterClockwise
	px := -1
	if p.Arg.TargetType == damage.TargetPlayer {
		// 敵の攻撃
		sx = battlecommon.FieldNum.X - 2
		sy = p.Arg.GetObjectPos(p.Arg.OwnerID).Y
		if sy == battlecommon.FieldNum.Y-1 {
			act = boomerangActTypeClockwise
		}
		px = battlecommon.FieldNum.X - 1
	}

	p.actType = act
	p.pos = point.Point{X: sx, Y: sy}
	p.next = point.Point{X: sx, Y: sy}
	p.prev = point.Point{X: px, Y: sy}
}

func (p *Boomerang) Process() (bool, error) {
	if p.count == 0 {
		p.Arg.SoundOn(resources.SEBoomerangThrow)
	}

	if p.count%boomerangNextStepCount == 0 {
		// Update current pos
		p.prev = p.pos
		p.pos = p.next

		p.Arg.DamageMgr.New(damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           boomerangNextStepCount + 1,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeSpreadHit,
			ShowHitArea:   false,
			Element:       damage.ElementWood,
		})

		switch p.actType {
		case boomerangActTypeCounterClockwise:
			if p.next.Y == 0 {
				if p.next.X == 0 && p.turnNum < 2 {
					p.turnNum++
					p.next.Y++
				} else {
					if p.next.X == battlecommon.FieldNum.X-1 {
						p.turnNum++
					}
					p.next.X--
				}
			} else if p.next.Y == battlecommon.FieldNum.Y-1 {
				if p.next.X == battlecommon.FieldNum.X-1 && p.turnNum < 2 {
					p.turnNum++
					p.next.Y--
				} else {
					if p.next.X == 0 {
						p.turnNum++
					}
					p.next.X++
				}
			} else {
				if p.next.X == 0 {
					p.next.Y++
				} else {
					p.next.Y--
				}
			}
		case boomerangActTypeClockwise:
			if p.next.Y == 0 {
				if p.next.X == battlecommon.FieldNum.X-1 && p.turnNum < 2 {
					p.turnNum++
					p.next.Y++
				} else {
					if p.next.X == 0 {
						p.turnNum++
					}
					p.next.X++
				}
			} else if p.next.Y == battlecommon.FieldNum.Y-1 {
				if p.next.X == 0 && p.turnNum < 2 {
					p.turnNum++
					p.next.Y--
				} else {
					if p.next.X == battlecommon.FieldNum.X-1 {
						p.turnNum++
					}
					p.next.X--
				}
			} else {
				if p.next.X == 0 {
					p.next.Y--
				} else {
					p.next.Y++
				}
			}
		default:
			system.SetError("not implemented yet")
		}
	}

	p.count++
	if p.pos.X < 0 || p.pos.X >= battlecommon.FieldNum.X || p.pos.Y < 0 || p.pos.Y >= battlecommon.FieldNum.Y {
		return true, nil
	}
	return false, nil
}

func (p *Boomerang) GetCount() int {
	return p.count
}

func (p *Boomerang) GetPos() (prev, current, next point.Point) {
	return p.prev, p.pos, p.next
}

func (p *Boomerang) GetNextStepCount() int {
	return boomerangNextStepCount
}
