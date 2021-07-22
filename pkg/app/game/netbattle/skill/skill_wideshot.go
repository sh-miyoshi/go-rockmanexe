package skill

import (
	"github.com/google/uuid"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type wideShot struct {
	bodyID        string
	beginID       string
	x             int
	y             int
	count         int
	power         int
	nextStepCount int
}

func newWideShot(x, y int, power int, nextStepCount int) *wideShot {
	return &wideShot{
		bodyID:        uuid.New().String(),
		beginID:       uuid.New().String(),
		x:             x,
		y:             y,
		power:         power,
		nextStepCount: nextStepCount,
	}
}

func (p *wideShot) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SEWideShot)
		// Add wide shot body
		netconn.SendObject(object.Object{
			ID:             p.bodyID,
			Type:           object.TypeWideShotBody,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       40,
			ViewOfsY:       -13,
		})

		// Add wide shot begin
		netconn.SendObject(object.Object{
			ID:             p.beginID,
			Type:           object.TypeWideShotBegin,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       62,
			ViewOfsY:       20,
		})
	}

	num, delay := draw.GetImageInfo(object.TypeWideShotBody)

	if p.count == num*delay {
		netconn.RemoveObject(p.bodyID)
		netconn.RemoveObject(p.beginID)
	}

	// Wide Shot Move
	if p.count > num*delay {
		if p.count%p.nextStepCount == 0 {
			p.x++
			if p.x >= appfield.FieldNumX {
				return true, nil
			}

			// Add object
			netconn.SendObject(object.Object{
				ID:             p.beginID,
				Type:           object.TypeWideShotMove,
				X:              p.x,
				Y:              p.y,
				UpdateBaseTime: true,
				ViewOfsY:       20,
				Speed:          p.nextStepCount,
			})

			// TODO Add damage
		}
	}

	return false, nil
}

func (p *wideShot) RemoveObject() {
	netconn.RemoveObject(p.bodyID)
	netconn.RemoveObject(p.beginID)
}
