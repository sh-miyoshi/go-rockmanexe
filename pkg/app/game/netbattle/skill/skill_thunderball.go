package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

const (
	nextStepCount = 80
	maxMoveCount  = 6
)

type thunderBall struct {
	id      string
	x       int
	y       int
	count   int
	power   int
	moveCnt int
	nextX   int
	nextY   int
}

func newThunderBall(x, y int, power int) *thunderBall {
	return &thunderBall{
		id:    uuid.New().String(),
		x:     x,
		y:     y,
		power: power,
		nextX: x + 1,
		nextY: y,
	}
}

func (p *thunderBall) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SEThunderBall)

		// Add object
		netconn.SendObject(p.getObject())
	}

	halfNext := nextStepCount / 2
	if p.count%nextStepCount == 0 {
		p.x = p.nextX
		p.y = p.nextY

		objs := getEnemies()
		if len(objs) == 0 {
			// no target
			p.nextX++
		} else {
			xdif := objs[0].X - p.x
			ydif := objs[0].Y - p.y

			if xdif != 0 || ydif != 0 {
				if common.Abs(xdif) > common.Abs(ydif) {
					// move to x
					p.nextX += (xdif / common.Abs(xdif))
				} else {
					// move to y
					p.nextY += (ydif / common.Abs(ydif))
				}
			}
		}

		p.moveCnt++
		if p.x < 0 || p.x >= field.FieldNumX || p.y < 0 || p.y > field.FieldNumY || p.moveCnt > maxMoveCount {
			return true, nil
		}

		netconn.SendObject(p.getObject())
	}

	if p.count%halfNext == 0 {
		x := p.x
		y := p.y
		if p.count%nextStepCount >= halfNext {
			x = p.nextX
			y = p.nextY
		}

		netconn.SendDamages([]damage.Damage{
			{
				ID:          uuid.New().String(),
				PosX:        x,
				PosY:        y,
				Power:       p.power,
				TTL:         halfNext,
				TargetType:  damage.TargetOtherClient,
				ShowHitArea: true,
			},
		})

		if isObjectHit(x, y) {
			return true, nil
		}
	}

	p.count++
	return false, nil
}

func (p *thunderBall) RemoveObject() {
	netconn.RemoveObject(p.id)
}

func (p *thunderBall) StopByPlayer() {
}

func (p *thunderBall) getObject() object.Object {
	return object.Object{
		ID:             p.id,
		Type:           object.TypeThunderBall,
		X:              p.x,
		Y:              p.y,
		UpdateBaseTime: true,
		TargetX:        p.nextX,
		TargetY:        p.nextY,
		Speed:          nextStepCount,
	}
}
