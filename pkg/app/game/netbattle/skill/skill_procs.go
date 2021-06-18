package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type cannon struct {
	id       string
	x        int
	y        int
	count    int
	timeInit bool
}

func newCannon(id string, x, y int) *cannon {
	return &cannon{
		id:       id,
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
			ID:             p.id,
			Type:           field.ObjectTypeCannonAtk,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: p.timeInit,
		},
		// Body
		{
			ID:             p.id,
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
