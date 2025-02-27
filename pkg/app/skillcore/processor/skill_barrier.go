package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type Barrier struct {
	Arg skillcore.Argument

	count int
	name  string
	power int
}

func (p *Barrier) Init(skillID int) {
	switch skillID {
	case resources.SkillBarrier:
		p.name = "バリア"
		p.power = 10
	case resources.SkillBarrier100:
		p.name = "バリア100"
		p.power = 100
	case resources.SkillBarrier200:
		p.name = "バリア200"
		p.power = 200
	}
}

func (p *Barrier) Update() (bool, error) {
	p.count++

	showTm := 60
	if p.count == 1 {
		p.Arg.Cutin(p.name, showTm)
		// WIP: add barrier
	}

	return p.count > showTm, nil
}

func (p *Barrier) GetCount() int {
	return p.count
}
