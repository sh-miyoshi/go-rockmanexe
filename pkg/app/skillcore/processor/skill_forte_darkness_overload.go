package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

const (
	forteDarknessOverloadEndCount = 64
)

type ForteDarknessOverload struct {
	Arg skillcore.Argument

	count int
}

func (p *ForteDarknessOverload) Process() (bool, error) {
	p.count++
	if p.count >= forteDarknessOverloadEndCount {
		// WIP: damage
		return true, nil
	}

	return false, nil
}

func (p *ForteDarknessOverload) GetCount() int {
	return p.count
}
