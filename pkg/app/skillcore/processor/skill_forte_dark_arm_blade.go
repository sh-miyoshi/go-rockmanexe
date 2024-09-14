package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteDarkArmBladeEndCount = 16
)

type ForteDarkArmBlade struct {
	Arg skillcore.Argument

	count  int
	atkPos point.Point
}

func (p *ForteDarkArmBlade) Init(skillID int) {
	p.count = 0
	p.atkPos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	switch skillID {
	case resources.SkillForteDarkArmBladeType1:
		p.atkPos.X--
	case resources.SkillForteDarkArmBladeType2:
	}
}

func (p *ForteDarkArmBlade) Process() (bool, error) {
	p.count++
	return p.count >= forteDarkArmBladeEndCount, nil
}

func (p *ForteDarkArmBlade) GetCount() int {
	return p.count
}

func (p *ForteDarkArmBlade) GetPos() point.Point {
	return p.atkPos
}
