package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type Invisible struct {
	Arg skillcore.Argument

	count int
}

func (p *Invisible) Process() (bool, error) {
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

func (p *Invisible) GetEndCount() int {
	// EndCount = delay * (len(img) + keepCount)
	// TODO
	return 1
}
