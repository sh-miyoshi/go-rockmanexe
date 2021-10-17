package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type longSword struct {
	id    string
	x     int
	y     int
	count int
	power int
}

func newLongSword(x, y int, power int) *longSword {
	return &longSword{
		id:    uuid.New().String(),
		x:     x,
		y:     y,
		power: power,
	}
}

func (p *longSword) Process() (bool, error) {
	p.count++

	if p.count == 5 {
		// Add object
		netconn.SendObject(object.Object{
			ID:             p.id,
			Type:           object.TypeLongSword,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       100,
		})
	}

	// num and delay are the same for normal, wide, and long longSword
	delay := object.ImageDelays[object.TypeLongSword]
	num := 4

	if p.count == 1*delay {
		p.addDamage()
	}

	if p.count > num*delay {
		return true, nil
	}
	return false, nil
}

func (p *longSword) GetObjects() []object.Object {
	return []object.Object{
		{
			ID:             p.id,
			Type:           object.TypeLongSword,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       100,
		},
	}
}

func (p *longSword) addDamage() {
	damages := []damage.Damage{}

	dm := damage.Damage{
		ID:         uuid.New().String(),
		Power:      p.power,
		TTL:        1,
		TargetType: damage.TargetOtherClient,
		BigDamage:  true,
	}

	dm.PosX = p.x + 1
	dm.PosY = p.y
	damages = append(damages, dm)
	dm.PosX = p.x + 2
	damages = append(damages, dm)

	netconn.SendDamages(damages)
}
