package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	hitNum      = 8
	atkInterval = 4
)

type Tornado struct {
	SkillID int
	Arg     skillcore.Argument

	count     int
	atkCount  int
	objPos    point.Point
	targetPos point.Point
}

func (p *Tornado) Init() {
	p.objPos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	p.targetPos = point.Point{X: p.objPos.X + 2, Y: p.objPos.Y}
}

func (p *Tornado) Process() (bool, error) {
	p.count++

	if p.count == 1 {
		p.Arg.SoundOn(resources.SETornado)
	}

	if p.count%atkInterval == 0 {
		lastAtk := p.atkCount == hitNum-1
		p.Arg.DamageMgr.New(damage.Damage{
			DamageType:    damage.TypePosition,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			BigDamage:     lastAtk,
			Element:       damage.ElementNone,
			Pos:           p.targetPos,
			TTL:           atkInterval,
			ShowHitArea:   false,
		})

		p.atkCount++
		return p.atkCount >= hitNum, nil
	}

	return false, nil
}

func (p *Tornado) GetCount() int {
	return p.count
}

func (p *Tornado) GetEndCount() int {
	return 32
}

func (p *Tornado) GetPos() (obj, target point.Point) {
	return p.objPos, p.targetPos
}
