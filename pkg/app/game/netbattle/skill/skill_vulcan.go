package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/netconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
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

			// TODO add damage

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
