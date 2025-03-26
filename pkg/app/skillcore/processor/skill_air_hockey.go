package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	airHockeyNextStepCount = 80
)

type AirHockey struct {
	SkillID int
	Arg     skillcore.Argument

	count int
	pos   point.Point
	next  point.Point
	prev  point.Point
}

func (p *AirHockey) Init() {
	pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
	x := pos.X + 1
	if p.Arg.TargetType == damage.TargetPlayer {
		x = pos.X - 1
	}

	first := point.Point{X: x, Y: pos.Y}
	p.pos = first
	p.prev = pos
	p.next = first
}

func (p *AirHockey) Update() (bool, error) {
	p.count++
	return false, nil
}

func (p *AirHockey) GetCount() int {
	return p.count
}

func (p *AirHockey) GetPos() (prev, current, next point.Point) {
	return p.prev, p.pos, p.next
}

func (p *AirHockey) GetNextStepCount() int {
	return airHockeyNextStepCount
}
