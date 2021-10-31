package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
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
	// TODO(敵の場合, 直線移動)
	sx := 0
	sy := field.FieldNumY - 1
	act := boomerangActTypeCounterClockwise
	px := -1

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

	sx := field.PanelSizeX*cnt/boomerangNextStepCount - field.PanelSizeX/2
	sy := field.PanelSizeY*cnt/boomerangNextStepCount - field.PanelSizeY/2

	var ofsx, ofsy int
	if cnt < boomerangNextStepCount/2 {
		ofsx = sx * (p.x - p.prevX)
		ofsy = sy * (p.y - p.prevY)
	} else {
		ofsx = (sx) * (p.nextX - p.x)
		ofsy = (sy) * (p.nextY - p.y)
	}

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

		// TODO damage

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
