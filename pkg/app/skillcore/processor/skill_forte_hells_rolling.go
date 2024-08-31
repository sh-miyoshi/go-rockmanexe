package processor

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"

type ForteHellsRolling struct {
	Arg skillcore.Argument

	count int
}

func (p *ForteHellsRolling) Process() (bool, error) {
	p.count++
	return false, nil
}

func (p *ForteHellsRolling) GetCount() int {
	return p.count
}
