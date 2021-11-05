package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type boomerang struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int
	ActType    int

	count   int
	turnNum int
	x       int
	y       int
	nextX   int
	nextY   int
	prevX   int
	prevY   int
}

const (
	boomerangNextStepCount = 6
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
	sy := field.FieldNumY - 1
	act := boomerangActTypeCounterClockwise
	px := -1
	if arg.TargetType == damage.TargetPlayer {
		// 敵の攻撃
		sx = field.FieldNumX - 2
		_, sy = objanim.GetObjPos(arg.OwnerID)
		if sy == field.FieldNumY-1 {
			act = boomerangActTypeClockwise
		}
		px = field.FieldNumX - 1
	}

	return &boomerang{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
		ActType:    act,
		count:      0,
		turnNum:    0,
		x:          sx,
		y:          sy,
		nextX:      sx,
		nextY:      sy,
		prevX:      px,
		prevY:      sy,
	}
}

func (p *boomerang) Draw() {
	x, y := battlecommon.ViewPos(p.x, p.y)
	n := (p.count / delayBoomerang) % len(imgBoomerang)

	cnt := p.count % boomerangNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(p.nextX, p.x, p.prevX, cnt, boomerangNextStepCount, field.PanelSizeX)
	ofsy := battlecommon.GetOffset(p.nextY, p.y, p.prevY, cnt, boomerangNextStepCount, field.PanelSizeY)
	dxlib.DrawRotaGraph(x+int32(ofsx), y+25+int32(ofsy), 1, 0, imgBoomerang[n], dxlib.TRUE)
}

func (p *boomerang) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SEBoomerangThrow)
	}

	if p.count%boomerangNextStepCount == 0 {
		// Update current pos
		p.prevX = p.x
		p.prevY = p.y
		p.x = p.nextX
		p.y = p.nextY

		damage.New(damage.Damage{
			PosX:          p.x,
			PosY:          p.y,
			Power:         int(p.Power),
			TTL:           boomerangNextStepCount + 1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeSpreadHit,
			ShowHitArea:   false,
		})

		switch p.ActType {
		case boomerangActTypeCounterClockwise:
			if p.nextY == 0 {
				if p.nextX == 0 && p.turnNum < 2 {
					p.turnNum++
					p.nextY++
				} else {
					if p.nextX == field.FieldNumX-1 {
						p.turnNum++
					}
					p.nextX--
				}
			} else if p.nextY == field.FieldNumY-1 {
				if p.nextX == field.FieldNumX-1 && p.turnNum < 2 {
					p.turnNum++
					p.nextY--
				} else {
					if p.nextX == 0 {
						p.turnNum++
					}
					p.nextX++
				}
			} else {
				if p.nextX == 0 {
					p.nextY++
				} else {
					p.nextY--
				}
			}
		case boomerangActTypeClockwise:
			if p.nextY == 0 {
				if p.nextX == field.FieldNumX-1 && p.turnNum < 2 {
					p.turnNum++
					p.nextY++
				} else {
					if p.nextX == 0 {
						p.turnNum++
					}
					p.nextX++
				}
			} else if p.nextY == field.FieldNumY-1 {
				if p.nextX == 0 && p.turnNum < 2 {
					p.turnNum++
					p.nextY--
				} else {
					if p.nextX == field.FieldNumX-1 {
						p.turnNum++
					}
					p.nextX--
				}
			} else {
				if p.nextX == 0 {
					p.nextY--
				} else {
					p.nextY++
				}
			}
		default:
			panic("not implemented yet")
		}
	}

	p.count++
	if p.x < 0 || p.x >= field.FieldNumX || p.y < 0 || p.y >= field.FieldNumY {
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
