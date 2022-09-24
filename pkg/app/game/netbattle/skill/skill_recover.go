package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type recover struct {
	id    string
	x     int
	y     int
	count int
	power int
}

func newRecover(x, y int, power int) *recover {
	return &recover{
		id:    uuid.New().String(),
		x:     x,
		y:     y,
		power: power,
	}
}

func (p *recover) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SERecover)

		// Add object
		netconn.GetInst().SendObject(object.Object{
			ID:             p.id,
			Type:           object.TypeRecover,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
		})

		// Add damage
		netconn.GetInst().AddDamage(damage.Damage{
			ID:         uuid.New().String(),
			PosX:       p.x,
			PosY:       p.y,
			Power:      -p.power,
			TTL:        1,
			TargetType: damage.TargetOwn,
		})
	}

	num, delay := draw.GetInst().GetObjectImageInfo(object.TypeRecover)
	if p.count > num*delay {
		return true, nil
	}
	return false, nil
}

func (p *recover) RemoveObject() {
	netconn.GetInst().RemoveObject(p.id)
}

func (p *recover) StopByPlayer() {
}
