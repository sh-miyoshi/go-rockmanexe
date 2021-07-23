package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/object"
)

type shockWave struct {
	pickID string
	waveID string
	x      int
	y      int
	count  int
	power  int
}

func newShockWave(x, y int, power int) *shockWave {
	return &shockWave{
		pickID: uuid.New().String(),
		waveID: uuid.New().String(),
		x:      x,
		y:      y,
		power:  power,
	}
}

func (p *shockWave) Process() (bool, error) {
	p.count++
	waveNum, waveDelay := draw.GetImageInfo(object.TypeShockWave)
	pickNum, pickDelay := draw.GetImageInfo(object.TypePick)

	if p.count == 1 {
		// Add pick
		netconn.SendObject(object.Object{
			ID:             p.pickID,
			Type:           object.TypePick,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsY:       -15,
		})
	}

	if p.count == pickNum*pickDelay+1 {
		netconn.RemoveObject(p.pickID)
	}

	if p.count > 10 {
		n := waveNum * waveDelay
		if p.count%(n) == 11 {
			p.x++
			if p.x >= appfield.FieldNumX {
				return true, nil
			}

			sound.On(sound.SEShockWave)
			// Add wave
			netconn.SendObject(object.Object{
				ID:             p.waveID,
				Type:           object.TypeShockWave,
				X:              p.x,
				Y:              p.y,
				UpdateBaseTime: true,
				ShowHitArea:    true,
			})

			// Add damage
			netconn.SendDamages([]damage.Damage{
				{
					ID:         uuid.New().String(),
					ClientID:   config.Get().Net.ClientID,
					PosX:       p.x,
					PosY:       p.y,
					Power:      p.power,
					TTL:        n - 2,
					TargetType: damage.TargetOtherClient,
				},
			})
		}
	}

	return false, nil
}

func (p *shockWave) RemoveObject() {
	netconn.RemoveObject(p.pickID)
	netconn.RemoveObject(p.waveID)
}
