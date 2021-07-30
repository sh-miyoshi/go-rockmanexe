package skill

import (
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

type spreadGun struct {
	atkID  string
	bodyID string
	x      int
	y      int
	count  int
	power  int
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

		for x := p.x + 1; x < appfield.FieldNumX; x++ {
			pn := netfield.GetPanelInfo(x, p.y)
			if pn.ObjectID != "" {
				// Hit
				sound.On(sound.SESpreadHit)

				damages := []damage.Damage{}
				dm := damage.Damage{
					ID:            uuid.New().String(),
					PosX:          x,
					PosY:          p.y,
					Power:         p.power,
					TTL:           1,
					TargetType:    damage.TargetOtherClient,
					HitEffectType: effect.TypeHitBigEffect,
				}
				damages = append(damages, dm)

				// Spreading
				dm.HitEffectType = 0
				for sy := -1; sy <= 1; sy++ {
					if p.y+sy < 0 || p.y+sy >= appfield.FieldNumY {
						continue
					}
					for sx := -1; sx <= 1; sx++ {
						if sy == 0 && sx == 0 {
							continue
						}
						if x+sx >= 0 && x+sx < appfield.FieldNumX {
							// Send effect
							netconn.SendEffect(effect.Effect{
								ID:   uuid.New().String(),
								Type: effect.TypeSpreadHitEffect,
								X:    x + sx,
								Y:    p.y + sy,
							})

							// Add damage
							dm.PosX = x + sx
							dm.PosY = p.y + sy
							damages = append(damages, dm)
						}
					}
				}

				netconn.SendDamages(damages)
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

	if p.count > max {
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
