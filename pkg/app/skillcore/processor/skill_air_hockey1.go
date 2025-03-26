package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type AirHockey1 struct {
	Arg skillcore.Argument

	count int
}

func (p *AirHockey1) Update() (bool, error) {
	p.count++
	return true, nil
}

func (p *AirHockey1) GetCount() int {
	return p.count
}
