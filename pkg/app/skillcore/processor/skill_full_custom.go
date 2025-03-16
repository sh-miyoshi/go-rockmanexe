package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type FullCustom struct {
	Arg skillcore.Argument

	count int
}

func (p *FullCustom) Update() (bool, error) {
	p.count++
	return true, nil
}

func (p *FullCustom) GetCount() int {
	return p.count
}
