package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/object"
)

type sword struct {
	id    string
	x     int
	y     int
	count int
	typ   int
	power int
}

func newSword(x, y int, power int, typ int) *sword {
	return &sword{
		id:    uuid.New().String(),
		x:     x,
		y:     y,
		typ:   typ,
		power: power,
	}
}

func (p *sword) Process() (bool, error) {
	p.count++

	if p.count == 5 {
		// Add object
		objType := -1
		switch p.typ {
		case skill.TypeSword:
			objType = object.TypeSword
		case skill.TypeWideSword:
			objType = object.TypeWideSword
		case skill.TypeLongSword:
			objType = object.TypeLongSword
		}
		net.GetInst().SendObject(object.Object{
			ID:             p.id,
			Type:           objType,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       100,
		})
	}

	// num and delay are the same for normal, wide, and long sword
	num, delay := draw.GetInst().GetObjectImageInfo(object.TypeSword)

	if p.count == 1*delay {
		sound.On(sound.SESword)
		p.addDamage()
	}

	if p.count > num*delay {
		return true, nil
	}
	return false, nil
}

func (p *sword) RemoveObject() {
	net.GetInst().RemoveObject(p.id)
}

func (p *sword) addDamage() {
	dm := damage.Damage{
		ID:         uuid.New().String(),
		Power:      p.power,
		TTL:        1,
		TargetType: damage.TargetOtherClient,
		BigDamage:  true,
	}

	dm.PosX = p.x + 1
	dm.PosY = p.y
	net.GetInst().AddDamage(dm)

	switch p.typ {
	case skill.TypeSword:
		// No more damage area
	case skill.TypeWideSword:
		dm.PosY = p.y - 1
		net.GetInst().AddDamage(dm)
		dm.PosY = p.y + 1
		net.GetInst().AddDamage(dm)
	case skill.TypeLongSword:
		dm.PosX = p.x + 2
		net.GetInst().AddDamage(dm)
	}
}

func (p *sword) StopByPlayer() {
	p.RemoveObject()
}
