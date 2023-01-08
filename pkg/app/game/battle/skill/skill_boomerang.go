package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type boomerang struct {
	ID      string
	Arg     Argument
	ActType int

	count   int
	turnNum int
	pos     common.Point
	next    common.Point
	prev    common.Point
}

const (
	boomerangNextStepCount = 6
	delayBoomerang         = 8
)

const (
	boomerangActTypeClockwise = iota
	boomerangActTypeCounterClockwise
	boomerangActTypeLeftStraight
	boomerangActTypeRightStraight
)

func newBoomerang(objID string, arg Argument) *boomerang {
	// TODO(直線移動)
	sx := 0
	sy := field.FieldNum.Y - 1
	act := boomerangActTypeCounterClockwise
	px := -1
	if arg.TargetType == damage.TargetPlayer {
		// 敵の攻撃
		sx = field.FieldNum.X - 2
		sy = objanim.GetObjPos(arg.OwnerID).Y
		if sy == field.FieldNum.Y-1 {
			act = boomerangActTypeClockwise
		}
		px = field.FieldNum.X - 1
	}

	return &boomerang{
		ID:      objID,
		Arg:     arg,
		ActType: act,
		count:   0,
		turnNum: 0,
		pos:     common.Point{X: sx, Y: sy},
		next:    common.Point{X: sx, Y: sy},
		prev:    common.Point{X: px, Y: sy},
	}
}

func (p *boomerang) Draw() {
	view := battlecommon.ViewPos(p.pos)
	n := (p.count / delayBoomerang) % len(imgBoomerang)

	cnt := p.count % boomerangNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(p.next.X, p.pos.X, p.prev.X, cnt, boomerangNextStepCount, field.PanelSize.X)
	ofsy := battlecommon.GetOffset(p.next.Y, p.pos.Y, p.prev.Y, cnt, boomerangNextStepCount, field.PanelSize.Y)
	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+25+ofsy, 1, 0, imgBoomerang[n], true)
}

func (p *boomerang) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SEBoomerangThrow)
	}

	if p.count%boomerangNextStepCount == 0 {
		// Update current pos
		p.prev = p.pos
		p.pos = p.next

		damage.New(damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           boomerangNextStepCount + 1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeSpreadHit,
			ShowHitArea:   false,
			DamageType:    damage.TypeWood,
		})

		switch p.ActType {
		case boomerangActTypeCounterClockwise:
			if p.next.Y == 0 {
				if p.next.X == 0 && p.turnNum < 2 {
					p.turnNum++
					p.next.Y++
				} else {
					if p.next.X == field.FieldNum.X-1 {
						p.turnNum++
					}
					p.next.X--
				}
			} else if p.next.Y == field.FieldNum.Y-1 {
				if p.next.X == field.FieldNum.X-1 && p.turnNum < 2 {
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
				if p.next.X == field.FieldNum.X-1 && p.turnNum < 2 {
					p.turnNum++
					p.next.Y++
				} else {
					if p.next.X == 0 {
						p.turnNum++
					}
					p.next.X++
				}
			} else if p.next.Y == field.FieldNum.Y-1 {
				if p.next.X == 0 && p.turnNum < 2 {
					p.turnNum++
					p.next.Y--
				} else {
					if p.next.X == field.FieldNum.X-1 {
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
	if p.pos.X < 0 || p.pos.X >= field.FieldNum.X || p.pos.Y < 0 || p.pos.Y >= field.FieldNum.Y {
		return true, nil
	}
	return false, nil
}

func (p *boomerang) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *boomerang) StopByOwner() {
	// Nothing to do after throwing
}
