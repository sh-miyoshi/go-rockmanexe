package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
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
		netconn.SendObject(object.Object{
			ID:             p.id,
			Type:           objType,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       100,
		})
	}

	// num and delay are the same for normal, wide, and long sword
	num, delay := draw.GetImageInfo(object.TypeSword)

	if p.count == 1*delay {
		sound.On(sound.SESword)

		// TODO add damage
	}

	if p.count > num*delay {
		return true, nil
	}
	return false, nil
}

func (p *sword) RemoveObject() {
	netconn.RemoveObject(p.id)
}
