package skill

import (
	"math/rand"

	"github.com/google/uuid"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type cannon struct {
	atkID  string
	bodyID string
	x      int
	y      int
	count  int
	typ    int
	power  int
}

func newCannon(x, y int, power int, typ int) *cannon {
	return &cannon{
		atkID:  uuid.New().String(),
		bodyID: uuid.New().String(),
		x:      x,
		y:      y,
		typ:    typ,
		power:  power,
	}
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		netconn.SendObject(p.getObjectInfo(true, false, true)) // Body
	}

	if p.count == 15 {
		netconn.SendObject(p.getObjectInfo(false, false, true)) // Attack
	}

	if p.count == 20 {
		sound.On(sound.SECannon)
		p.addDamage()
	}

	// num and delay are the same for normal, high, and mega
	bodyNum, bodyDelay := draw.GetImageInfo(object.TypeNormalCannonBody)
	atkNum, atkDelay := draw.GetImageInfo(object.TypeNormalCannonAtk)
	max := bodyNum * bodyDelay
	if n := atkNum*atkDelay + 15; max < n {
		max = n
	}

	if p.count == 2*bodyDelay {
		netconn.SendObject(p.getObjectInfo(true, true, false)) // Shifted Body
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

func (p *cannon) StopByPlayer() {
	p.RemoveObject()
}

func (p *cannon) addDamage() {
	dm := []damage.Damage{}
	for x := p.x + 1; x < appfield.FieldNumX; x++ {
		dm = append(dm, damage.Damage{
			ID:            uuid.New().String(),
			PosX:          x,
			PosY:          p.y,
			Power:         p.power,
			TTL:           1,
			TargetType:    damage.TargetOtherClient,
			HitEffectType: effect.TypeCannonHitEffect,
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

func (p *cannon) getObjectInfo(isBody bool, isShift bool, updateTime bool) object.Object {
	id := p.atkID
	vx := int32(90)
	vy := int32(-10)
	if isBody {
		id = p.bodyID
		if isShift {
			vx = 33
		} else {
			vx = 48
		}
		vy = -12
	}

	typ := -1
	switch p.typ {
	case skill.TypeNormalCannon:
		if isBody {
			typ = object.TypeNormalCannonBody
		} else {
			typ = object.TypeNormalCannonAtk
		}
	case skill.TypeHighCannon:
		if isBody {
			typ = object.TypeHighCannonBody
		} else {
			typ = object.TypeHighCannonAtk
		}
	case skill.TypeMegaCannon:
		if isBody {
			typ = object.TypeMegaCannonBody
		} else {
			typ = object.TypeMegaCannonAtk
		}
	}

	return object.Object{
		ID:             id,
		Type:           typ,
		HP:             0,
		X:              p.x,
		Y:              p.y,
		ViewOfsX:       vx,
		ViewOfsY:       vy,
		UpdateBaseTime: updateTime,
	}
}
