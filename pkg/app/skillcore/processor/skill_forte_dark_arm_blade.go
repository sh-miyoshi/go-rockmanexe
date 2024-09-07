package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type ForteDarkArmBlade struct {
	Arg skillcore.Argument

	count int
}

func (p *ForteDarkArmBlade) Init(skillID int) {
	p.count = 0
}

func (p *ForteDarkArmBlade) Process() (bool, error) {
	p.count++
	return false, nil
}

func (p *ForteDarkArmBlade) GetCount() int {
	return p.count
}
