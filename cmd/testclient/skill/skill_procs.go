package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/cmd/testclient/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type cannon struct {
	atkID        string
	bodyID       string
	x            int
	y            int
	count        int
	timeInit     bool
	viewBodyOfsX int32
}

func newCannon(x, y int) *cannon {
	return &cannon{
		atkID:        uuid.New().String(),
		bodyID:       uuid.New().String(),
		x:            x,
		y:            y,
		timeInit:     true,
		viewBodyOfsX: 48,
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++

	// if p.count == 20 {
	// 	// TODO add damage
	// }

	bodyDelay := object.ImageDelays[object.TypeNormalCannonBody]
	atkDelay := object.ImageDelays[object.TypeNormalCannonAtk]
	bodyNum := 4
	atkNum := 8
	max := bodyNum * bodyDelay
	if n := atkNum*atkDelay + 15; max < n {
		max = n
	}

	if p.count == 2*bodyDelay {
		p.viewBodyOfsX = 33
		netconn.SendObject(object.Object{
			ID:             p.bodyID,
			Type:           object.TypeNormalCannonBody,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: p.timeInit,
			ViewOfsX:       p.viewBodyOfsX,
			ViewOfsY:       -12,
		})
	}

	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *cannon) GetObjects() []object.Object {
	res := []object.Object{
		// Attack
		{
			ID:             p.atkID,
			Type:           object.TypeNormalCannonAtk,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: p.timeInit,
			ViewOfsX:       90,
			ViewOfsY:       -10,
		},
		// Body
		{
			ID:             p.bodyID,
			Type:           object.TypeNormalCannonBody,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: p.timeInit,
			ViewOfsX:       p.viewBodyOfsX,
			ViewOfsY:       -12,
		},
	}

	p.timeInit = false

	return res
}
