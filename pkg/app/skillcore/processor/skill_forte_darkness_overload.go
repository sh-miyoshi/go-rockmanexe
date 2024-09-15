package processor

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"

type ForteDarknessOverload struct {
	Arg skillcore.Argument

	count int
}

func (p *ForteDarknessOverload) Process() (bool, error) {
	p.count++
	return false, nil
}

func (p *ForteDarknessOverload) GetCount() int {
	return p.count
}
