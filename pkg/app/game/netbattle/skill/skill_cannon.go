package skill

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/object"
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
		net.GetInst().SendObject(p.getObjectInfo(true, false, true)) // Body
	}

	if p.count == 15 {
		net.GetInst().SendObject(p.getObjectInfo(false, false, true)) // Attack
	}

	if p.count == 20 {
		sound.On(sound.SECannon)
		p.addDamage()
	}

	// num and delay are the same for normal, high, and mega
	bodyNum, bodyDelay := draw.GetInst().GetObjectImageInfo(object.TypeNormalCannonBody)
	atkNum, atkDelay := draw.GetInst().GetObjectImageInfo(object.TypeNormalCannonAtk)
	max := bodyNum * bodyDelay
	if n := atkNum*atkDelay + 15; max < n {
		max = n
	}

	if p.count == 2*bodyDelay {
		net.GetInst().SendObject(p.getObjectInfo(true, true, false)) // Shifted Body
	}

	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *cannon) RemoveObject() {
	net.GetInst().RemoveObject(p.atkID)
	net.GetInst().RemoveObject(p.bodyID)
}

func (p *cannon) StopByPlayer() {
	p.RemoveObject()
}

func (p *cannon) addDamage() {
	for x := p.x + 1; x < battlecommon.FieldNum.X; x++ {
		dm := damage.Damage{
			ID:            uuid.New().String(),
			PosX:          x,
			PosY:          p.y,
			Power:         p.power,
			TTL:           1,
			TargetType:    damage.TargetOtherClient,
			HitEffectType: effect.TypeCannonHitEffect,
			ViewOfsX:      rand.Intn(2*5) - 5,
			ViewOfsY:      rand.Intn(2*5) - 5,
			BigDamage:     true,
		}
		net.GetInst().AddDamage(dm)

		// break if object exists
		pn := netfield.GetPanelInfo(common.Point{X: x, Y: p.y})
		if pn.ObjectID != "" {
			sound.On(sound.SECannonHit)
			break
		}
	}
}

func (p *cannon) getObjectInfo(isBody bool, isShift bool, updateTime bool) object.Object {
	id := p.atkID
	vx := 90
	vy := -10
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
