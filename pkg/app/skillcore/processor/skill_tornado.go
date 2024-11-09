package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	tornadoHitNum      = 8
	tornadoAtkInterval = 4
)

type Tornado struct {
	Arg skillcore.Argument

	count     int
	atkCount  int
	objPos    point.Point
	targetPos point.Point
}

func (p *Tornado) Init() {
	p.objPos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	p.targetPos = point.Point{X: p.objPos.X + 2, Y: p.objPos.Y}
}

func (p *Tornado) Update() (bool, error) {
	p.count++

	if p.count == 1 {
		p.Arg.SoundOn(resources.SETornado)
	}

	if p.count%tornadoAtkInterval == 0 {
		lastAtk := p.atkCount == tornadoHitNum-1
		strengthType := damage.StrengthNone
		if lastAtk {
			strengthType = damage.StrengthHigh
		}

		p.Arg.DamageMgr.New(damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypePosition,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			StrengthType:  strengthType,
			Element:       damage.ElementNone,
			Pos:           p.targetPos,
			TTL:           tornadoAtkInterval,
			ShowHitArea:   false,
		})

		p.atkCount++
		return p.atkCount >= tornadoHitNum, nil
	}

	return false, nil
}

func (p *Tornado) GetCount() int {
	return p.count
}

func (p *Tornado) GetPos() (obj, target point.Point) {
	return p.objPos, p.targetPos
}
