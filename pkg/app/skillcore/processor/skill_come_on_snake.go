package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type ComeOnSnake struct {
	Arg skillcore.Argument

	count int
}

func (p *ComeOnSnake) Update() (bool, error) {
	p.count++
	return true, nil
}

func (p *ComeOnSnake) GetCount() int {
	return p.count
}
