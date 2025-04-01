package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type AirShoot struct {
	Arg skillcore.Argument

	count int
}

func (p *AirShoot) Update() (bool, error) {
	p.count++
	return true, nil
}

func (p *AirShoot) GetCount() int {
	return p.count
}
