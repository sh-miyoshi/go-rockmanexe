package skill

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/field"
)

type cannon struct {
	atkID  string
	bodyID string
	x      int
	y      int
	count  int
	power  int
}

func newCannon(x, y int, power int) *cannon {
	return &cannon{
		atkID:  uuid.New().String(),
		bodyID: uuid.New().String(),
		x:      x,
		y:      y,
		power:  power,
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		// Body
		netconn.SendObject(field.Object{
			ID:             p.bodyID,
			Type:           field.ObjectTypeCannonBody,
			HP:             0,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       48,
			ViewOfsY:       -12,
		})
	}

	if p.count == 15 {
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
	}

	if p.count == 20 {
		sound.On(sound.SECannon)
		p.addDamage()
	}

	bodyNum, bodyDelay := draw.GetImageInfo(field.ObjectTypeCannonBody)
	atkNum, atkDelay := draw.GetImageInfo(field.ObjectTypeCannonAtk)
	max := bodyNum * bodyDelay
	if n := atkNum*atkDelay + 15; max < n {
		max = n
	}

	if p.count == 2*bodyDelay {
		netconn.SendObject(field.Object{
			ID:       p.bodyID,
			Type:     field.ObjectTypeCannonBody,
			HP:       0,
			X:        p.x,
			Y:        p.y,
			ViewOfsX: 33,
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

func (p *cannon) addDamage() {
	clientID := config.Get().Net.ClientID

	dm := []damage.Damage{}
	for x := p.x + 1; x < appfield.FieldNumX; x++ {
		dm = append(dm, damage.Damage{
			ID:            uuid.New().String(),
			ClientID:      clientID,
			PosX:          x,
			PosY:          p.y,
			Power:         p.power,
			TTL:           1,
			TargetType:    damage.TargetOtherClient,
			HitEffectType: field.ObjectTypeCannonHitEffect,
			ViewOfsX:      int32(rand.Intn(2*5) - 5),
			ViewOfsY:      int32(rand.Intn(2*5) - 5),
		})

		// break if object exists
		pn := netfield.GetPanelInfo(x, p.y)
		if pn.ObjectID != "" {
			break
		}
	}
	netconn.SendDamages(dm)
}
