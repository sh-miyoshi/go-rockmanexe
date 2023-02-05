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

type wideShot struct {
	bodyID        string
	beginID       string
	moveID        string
	x             int
	y             int
	count         int
	power         int
	nextStepCount int
}

func newWideShot(x, y int, power int, nextStepCount int) *wideShot {
	return &wideShot{
		bodyID:        uuid.New().String(),
		beginID:       uuid.New().String(),
		moveID:        uuid.New().String(),
		x:             x,
		y:             y,
		power:         power,
		nextStepCount: nextStepCount,
	}
}

func (p *wideShot) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		sound.On(sound.SEWideShot)
		// Add wide shot body
		net.GetInst().SendObject(object.Object{
			ID:             p.bodyID,
			Type:           object.TypeWideShotBody,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       40,
			ViewOfsY:       -13,
		})

		// Add wide shot begin
		net.GetInst().SendObject(object.Object{
			ID:             p.beginID,
			Type:           object.TypeWideShotBegin,
			X:              p.x,
			Y:              p.y,
			UpdateBaseTime: true,
			ViewOfsX:       62,
			ViewOfsY:       20,
		})
	}

	num, delay := draw.GetInst().GetObjectImageInfo(object.TypeWideShotBody)

	if p.count == num*delay {
		net.GetInst().RemoveObject(p.bodyID)
		net.GetInst().RemoveObject(p.beginID)
	}

	// Wide Shot Move
	if p.count > num*delay {
		if p.count%p.nextStepCount == 0 {
			p.x++
			if p.x >= appfield.FieldNum.X {
				return true, nil
			}

			// Add object
			net.GetInst().SendObject(object.Object{
				ID:             p.moveID,
				Type:           object.TypeWideShotMove,
				X:              p.x,
				Y:              p.y,
				UpdateBaseTime: true,
				ViewOfsY:       20,
				Speed:          p.nextStepCount,
			})

			p.addDamages()

			if isObjectHit(p.x, p.y) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (p *wideShot) RemoveObject() {
	net.GetInst().RemoveObject(p.bodyID)
	net.GetInst().RemoveObject(p.beginID)
	net.GetInst().RemoveObject(p.moveID)
}

func (p *wideShot) StopByPlayer() {
	num, delay := draw.GetInst().GetObjectImageInfo(object.TypeWideShotBody)

	if p.count < num*delay {
		p.RemoveObject()
	}
}

func (p *wideShot) addDamages() {
	dm := damage.Damage{
		ID:         uuid.New().String(),
		PosX:       p.x,
		PosY:       p.y,
		Power:      p.power,
		TTL:        p.nextStepCount,
		TargetType: damage.TargetOtherClient,
		BigDamage:  true,
	}

	// Add damages to 3 wide
	net.GetInst().AddDamage(dm)
	dm.PosY--
	net.GetInst().AddDamage(dm)
	dm.PosY += 2
	net.GetInst().AddDamage(dm)
}
