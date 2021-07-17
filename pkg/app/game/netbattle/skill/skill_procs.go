package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type cannon struct {
	atkID        string
	bodyID       string
	x            int
	y            int
	count        int
	viewBodyOfsX int32
}

func newCannon(x, y int) *cannon {
	return &cannon{
		atkID:        uuid.New().String(),
		bodyID:       uuid.New().String(),
		x:            x,
		y:            y,
		viewBodyOfsX: 48,
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		// Attack
		netconn.SendObject(field.Object{
			ID:             p.atkID,
			Type:           field.ObjectTypeCannonAtk,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       90,
			ViewOfsY:       -10,
		})
		// Body
		netconn.SendObject(field.Object{
			ID:             p.bodyID,
			Type:           field.ObjectTypeCannonBody,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       p.viewBodyOfsX,
			ViewOfsY:       -12,
		})
	}

	if p.count == 20 {
		sound.On(sound.SECannon)
		// TODO add damage
	}

	bodyNum, bodyDelay := draw.GetImageInfo(field.ObjectTypeCannonBody)
	atkNum, atkDelay := draw.GetImageInfo(field.ObjectTypeCannonAtk)
	max := bodyNum * bodyDelay
	if n := atkNum*atkDelay + 15; max < n {
		max = n
	}

	if p.count == 2*bodyDelay {
		p.viewBodyOfsX = 33
		netconn.SendObject(field.Object{
			ID:       p.bodyID,
			Type:     field.ObjectTypeCannonBody,
			HP:       0,
			X:        p.x,
			Y:        p.y,
			ViewOfsX: p.viewBodyOfsX,
			ViewOfsY: -12,
		})
	}

	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *cannon) RemoveObject() {
	netconn.RemoveObject(p.atkID)
	netconn.RemoveObject(p.bodyID)
}
