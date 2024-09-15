package processor

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"

type ForteShootingBuster struct {
	Arg skillcore.Argument

	count int
}

func (p *ForteShootingBuster) Process() (bool, error) {
	p.count++
	return false, nil
}

func (p *ForteShootingBuster) GetCount() int {
	return p.count
}
