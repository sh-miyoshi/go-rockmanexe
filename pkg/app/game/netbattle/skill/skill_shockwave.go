package skill

import (
	"github.com/google/uuid"
	appfield "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/oldnet/object"
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
	waveNum, waveDelay := draw.GetInst().GetObjectImageInfo(object.TypeShockWave)
	pickNum, pickDelay := draw.GetInst().GetObjectImageInfo(object.TypePick)

	if p.count == 1 {
		// Add pick
		net.GetInst().SendObject(object.Object{
			ID:             p.pickID,
			Type:           object.TypePick,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsY:       -15,
		})
	}

	if p.count == pickNum*pickDelay+1 {
		net.GetInst().RemoveObject(p.pickID)
	}

	if p.count > 10 {
		n := waveNum * waveDelay
		if p.count%(n) == 11 {
			p.x++
			if p.x >= appfield.FieldNum.X {
				return true, nil
			}

			sound.On(sound.SEShockWave)
			// Add wave
			net.GetInst().SendObject(object.Object{
				ID:             p.waveID,
				Type:           object.TypeShockWave,
				X:              p.x,
				Y:              p.y,
				UpdateBaseTime: true,
			})

			// Add damage
			net.GetInst().AddDamage(damage.Damage{
				ID:          uuid.New().String(),
				PosX:        p.x,
				PosY:        p.y,
				Power:       p.power,
				TTL:         n - 2,
				TargetType:  damage.TargetOtherClient,
				ShowHitArea: true,
				BigDamage:   true,
			})
		}
	}

	return false, nil
}

func (p *shockWave) RemoveObject() {
	net.GetInst().RemoveObject(p.pickID)
	net.GetInst().RemoveObject(p.waveID)
}

func (p *shockWave) StopByPlayer() {
	if p.count < 10 {
		p.RemoveObject()
	}
}
