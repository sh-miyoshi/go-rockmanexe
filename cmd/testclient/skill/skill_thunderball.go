package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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
	prevX   int
	prevY   int
}

func newThunderBall(x, y int, power int) *thunderBall {
	return &thunderBall{
		id:    uuid.New().String(),
		x:     x + 1,
		y:     y,
		power: power,
		nextX: x + 1,
		nextY: y,
		prevX: x,
		prevY: y,
	}
}

func (p *thunderBall) Process() (bool, error) {
	// if p.count%nextStepCount == 2 {
	// 	if isObjectHit(p.x, p.y) {
	// 		return true, nil
	// 	}
	// }

	if p.count%nextStepCount == 0 {
		tx := p.x
		ty := p.y
		if p.count != 0 {
			// Update current pos
			p.prevX = p.x
			p.prevY = p.y
			p.x = p.nextX
			p.y = p.nextY

			p.moveCnt++
			if p.moveCnt > maxMoveCount {
				return true, nil
			}

			if p.x < 0 || p.x > 6 || p.y < 0 || p.y > 3 {
				return true, nil
			}
		}

		// Set next pos
		objs := getEnemies()
		if len(objs) == 0 {
			// no target
			p.nextX++
		} else {
			xdif := objs[0].X - tx
			ydif := objs[0].Y - ty

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

		netconn.SendDamages([]damage.Damage{
			{
				ID:          uuid.New().String(),
				PosX:        p.x,
				PosY:        p.y,
				Power:       p.power,
				TTL:         nextStepCount,
				TargetType:  damage.TargetOtherClient,
				ShowHitArea: true,
				BigDamage:   true, // TODO make paralysis
			},
		})

		netconn.SendObject(p.getObject())
	}

	p.count++
	return false, nil
}

func (p *thunderBall) GetObjects() []object.Object {
	return []object.Object{p.getObject()}
}

func (p *thunderBall) getObject() object.Object {
	logger.Debug("Object info: %+v", p)

	return object.Object{
		ID:             p.id,
		Type:           object.TypeThunderBall,
		X:              p.x,
		Y:              p.y,
		UpdateBaseTime: true,
		TargetX:        p.nextX,
		TargetY:        p.nextY,
		Speed:          nextStepCount,
		PrevX:          p.prevX,
		PrevY:          p.prevY,
	}
}