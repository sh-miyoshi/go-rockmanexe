package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

const (
	endCount = 60
)

type miniBomb struct {
	id    string
	x     int
	y     int
	count int
	power int
}

func newMiniBomb(x, y int, power int) *miniBomb {
	return &miniBomb{
		id:    uuid.New().String(),
		x:     x,
		y:     y,
		power: power,
	}
}

func (p *miniBomb) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		netconn.GetInst().SendObject(object.Object{
			ID:             p.id,
			Type:           object.TypeMiniBomb,
			X:              p.x,
			Y:              p.y,
			TargetX:        p.x + 3,
			TargetY:        p.y,
			UpdateBaseTime: true,
			Speed:          endCount,
		})
	}

	if p.count == endCount {
		// Add Explode
		sound.On(sound.SEExplode)

		netconn.GetInst().SendEffect(effect.Effect{
			ID:   uuid.New().String(),
			Type: effect.TypeExplodeEffect,
			X:    p.x + 3,
			Y:    p.y,
		})

		netconn.GetInst().AddDamage(damage.Damage{
			ID:         uuid.New().String(),
			PosX:       p.x + 3,
			PosY:       p.y,
			Power:      p.power,
			TTL:        1,
			TargetType: damage.TargetOtherClient,
			BigDamage:  true,
		})

		return true, nil
	}

	return false, nil
}

func (p *miniBomb) RemoveObject() {
	netconn.GetInst().RemoveObject(p.id)
}

func (p *miniBomb) StopByPlayer() {
	if p.count < 5 {
		p.RemoveObject()
	}
}
