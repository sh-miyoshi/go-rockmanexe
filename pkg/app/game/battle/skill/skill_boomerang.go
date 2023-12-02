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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type boomerang struct {
	ID      string
	Arg     Argument
	ActType int

	count   int
	turnNum int
	pos     point.Point
	next    point.Point
	prev    point.Point
	drawer  skilldraw.DrawBoomerang
}

const (
	boomerangActTypeClockwise = iota
	boomerangActTypeCounterClockwise
	boomerangActTypeLeftStraight
	boomerangActTypeRightStraight
)

func newBoomerang(objID string, arg Argument) *boomerang {
	// TODO(直線移動)
	sx := 0
	sy := battlecommon.FieldNum.Y - 1
	act := boomerangActTypeCounterClockwise
	px := -1
	if arg.TargetType == damage.TargetPlayer {
		// 敵の攻撃
		sx = battlecommon.FieldNum.X - 2
		sy = localanim.ObjAnimGetObjPos(arg.OwnerID).Y
		if sy == battlecommon.FieldNum.Y-1 {
			act = boomerangActTypeClockwise
		}
		px = battlecommon.FieldNum.X - 1
	}

	return &boomerang{
		ID:      objID,
		Arg:     arg,
		ActType: act,
		count:   0,
		turnNum: 0,
		pos:     point.Point{X: sx, Y: sy},
		next:    point.Point{X: sx, Y: sy},
		prev:    point.Point{X: px, Y: sy},
	}
}

func (p *boomerang) Draw() {
	p.drawer.Draw(p.prev, p.pos, p.next, p.count)
}

func (p *boomerang) Process() (bool, error) {
	if p.count == 0 {
		sound.On(resources.SEBoomerangThrow)
	}

	if p.count%resources.SkillBoomerangNextStepCount == 0 {
		// Update current pos
		p.prev = p.pos
		p.pos = p.next

		localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           resources.SkillBoomerangNextStepCount + 1,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeSpreadHit,
			ShowHitArea:   false,
			Element:       damage.ElementWood,
		})

		switch p.ActType {
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
			common.SetError("not implemented yet")
		}
	}

	p.count++
	if p.pos.X < 0 || p.pos.X >= battlecommon.FieldNum.X || p.pos.Y < 0 || p.pos.Y >= battlecommon.FieldNum.Y {
		return true, nil
	}
	return false, nil
}

func (p *boomerang) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *boomerang) StopByOwner() {
	// Nothing to do after throwing
}
