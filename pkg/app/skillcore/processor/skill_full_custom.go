package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type FullCustom struct {
	Arg skillcore.Argument
}

func (p *FullCustom) Update() (bool, error) {
	p.Arg.SetCustomGaugeMax(p.Arg.OwnerID)
	return true, nil
}

func (p *FullCustom) GetCount() int {
	return 0
}
