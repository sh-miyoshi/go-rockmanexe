package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	netfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

const (
	spreadWaitCount = 10
)

type spreadGun struct {
	atkID          string
	bodyID         string
	x              int
	y              int
	count          int
	power          int
	waitCount      int
	spreadBaseInfo damage.Damage
}

func newSpreadGun(x, y int, power int) *spreadGun {
	return &spreadGun{
		atkID:  uuid.New().String(),
		bodyID: uuid.New().String(),
		x:      x,
		y:      y,
		power:  power,
	}
}

func (p *spreadGun) Process() (bool, error) {
	p.count++

	if p.waitCount > 0 {
		p.waitCount++
		if p.waitCount >= spreadWaitCount {
			// Spreading
			damages := []damage.Damage{}
			dm := p.spreadBaseInfo
			dm.ID = uuid.New().String()
			dm.HitEffectType = 0
			x := dm.PosX
			y := dm.PosY
			for sy := -1; sy <= 1; sy++ {
				if y+sy < 0 || y+sy >= int(appfield.FieldNum.Y) {
					continue
				}
				for sx := -1; sx <= 1; sx++ {
					if sy == 0 && sx == 0 {
						continue
					}
					if x+sx >= 0 && x+sx < int(appfield.FieldNum.X) {
						// Send effect
						netconn.SendEffect(effect.Effect{
							ID:   uuid.New().String(),
							Type: effect.TypeSpreadHitEffect,
							X:    x + sx,
							Y:    y + sy,
						})

						// Add damage
						dm.PosX = x + sx
						dm.PosY = y + sy
						damages = append(damages, dm)
					}
				}
			}
			netconn.SendDamages(damages)
			p.waitCount = 0
		}
	}

	if p.count == 1 {
		// Add objects
		netconn.SendObject(object.Object{
			ID:             p.atkID,
			Type:           object.TypeSpreadGunAtk,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       100,
			ViewOfsY:       -20,
		})

		netconn.SendObject(object.Object{
			ID:             p.bodyID,
			Type:           object.TypeSpreadGunBody,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       50,
			ViewOfsY:       -18,
		})
	}

	if p.count == 5 {
		sound.On(sound.SEGun)

		for x := p.x + 1; x < int(appfield.FieldNum.X); x++ {
			pn := netfield.GetPanelInfo(common.Point{X: x, Y: p.y})
			if pn.ObjectID != "" {
				// Hit
				sound.On(sound.SESpreadHit)

				dm := damage.Damage{
					ID:            uuid.New().String(),
					PosX:          x,
					PosY:          p.y,
					Power:         p.power,
					TTL:           1,
					TargetType:    damage.TargetOtherClient,
					HitEffectType: effect.TypeHitBigEffect,
				}

				// Set spreading
				p.spreadBaseInfo = dm
				p.waitCount = 1

				netconn.SendDamages([]damage.Damage{dm})
				break
			}
		}
	}

	bodyNum, bodyDelay := draw.GetImageInfo(object.TypeSpreadGunBody)
	atkNum, atkDelay := draw.GetImageInfo(object.TypeSpreadGunAtk)
	max := bodyNum * bodyDelay
	if n := atkNum * atkDelay; max < n {
		max = n
	}

	if p.count > max && p.waitCount == 0 {
		return true, nil
	}
	return false, nil
}

func (p *spreadGun) RemoveObject() {
	netconn.RemoveObject(p.atkID)
	netconn.RemoveObject(p.bodyID)
}

func (p *spreadGun) StopByPlayer() {
	p.RemoveObject()
}
