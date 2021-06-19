package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type cannon struct {
	atkID    string
	bodyID   string
	x        int
	y        int
	count    int
	timeInit bool
}

func newCannon(x, y int) *cannon {
	return &cannon{
		atkID:    uuid.New().String(),
		bodyID:   uuid.New().String(),
		x:        x,
		y:        x,
		timeInit: true,
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++

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

	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *cannon) GetObjects() []field.Object {
	res := []field.Object{
		// Attack
		{
			ID:             p.atkID,
			Type:           field.ObjectTypeCannonAtk,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: p.timeInit,
		},
		// Body
		{
			ID:             p.bodyID,
			Type:           field.ObjectTypeCannonBody,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: p.timeInit,
		},
	}

	p.timeInit = false

	return res
}
