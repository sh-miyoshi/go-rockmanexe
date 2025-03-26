package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type AirHockey struct {
	Arg skillcore.Argument

	count int
}

func (p *AirHockey) Update() (bool, error) {
	p.count++
	return true, nil
}

func (p *AirHockey) GetCount() int {
	return p.count
}
