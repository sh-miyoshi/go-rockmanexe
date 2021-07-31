package skill

import (
	"math/rand"

	"github.com/google/uuid"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type vulcan struct {
	id       string
	x        int
	y        int
	count    int
	atkNum   int
	power    int
	atkCount int
	hit      bool
}

func newVulcan(x, y int, atkNum int) *vulcan {
	return &vulcan{
		id:     uuid.New().String(),
		x:      x,
		y:      y,
		power:  10,
		atkNum: atkNum,
	}
}

func (p *vulcan) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		// Body
		netconn.SendObject(object.Object{
			ID:             p.id,
			Type:           object.TypeVulcan,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
		})
	}

	_, delay := draw.GetImageInfo(object.TypeVulcan)
	if p.count >= delay*1 {
		if p.count%(delay*5) == delay*1 {
			sound.On(sound.SEGun)
			p.addDamage()

			p.atkCount++
			if p.atkCount == p.atkNum {
				return true, nil
			}
		}
	}

	return false, nil
}

func (p *vulcan) RemoveObject() {
	netconn.RemoveObject(p.id)
}

func (p *vulcan) StopByPlayer() {
	p.RemoveObject()
}

func (p *vulcan) addDamage() {
	dm := []damage.Damage{}
	hit := false
	eff := effect.Effect{}
	for x := p.x + 1; x < appfield.FieldNumX; x++ {
		pn := netfield.GetPanelInfo(x, p.y)
		if pn.ObjectID != "" {
			dm = append(dm, damage.Damage{
				ID:            uuid.New().String(),
				PosX:          x,
				PosY:          p.y,
				Power:         p.power,
				TTL:           1,
				TargetType:    damage.TargetOtherClient,
				HitEffectType: effect.TypeSpreadHitEffect,
			})
			eff = effect.Effect{
				ID:       uuid.New().String(),
				Type:     effect.TypeVulcanHit1Effect,
				X:        x,
				Y:        p.y,
				ViewOfsX: int32(rand.Intn(2*20) - 20),
				ViewOfsY: int32(rand.Intn(2*20) - 20),
			}
			if p.hit && x < appfield.FieldNumX-1 {
				dm = append(dm, damage.Damage{
					ID:            uuid.New().String(),
					PosX:          x + 1,
					PosY:          p.y,
					Power:         p.power,
					TTL:           1,
					TargetType:    damage.TargetOtherClient,
					HitEffectType: effect.TypeVulcanHit2Effect,
					ViewOfsX:      int32(rand.Intn(2*20) - 20),
					ViewOfsY:      int32(rand.Intn(2*20) - 20),
				})
			}
			hit = true
			sound.On(sound.SECannonHit)
			break
		}
	}
	p.hit = hit
	netconn.SendDamages(dm)
	if hit {
		netconn.SendEffect(eff)
	}
}
