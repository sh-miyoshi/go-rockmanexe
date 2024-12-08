package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type Invisible struct {
	Arg skillcore.Argument

	count int
}

func (p *Invisible) Update() (bool, error) {
	p.count++

	showTm := 60
	if p.count == 1 {
		p.Arg.Cutin("インビジブル", showTm)
		p.Arg.MakeInvisible(p.Arg.OwnerID, 6*60)
	}

	return p.count > showTm, nil
}

func (p *Invisible) GetCount() int {
	return p.count
}
