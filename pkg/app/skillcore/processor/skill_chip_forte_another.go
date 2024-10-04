package processor

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"

type ChipForteAnother struct {
	Arg skillcore.Argument

	count int
}

func (p *ChipForteAnother) Process() (bool, error) {
	return false, nil
}

func (p *ChipForteAnother) GetCount() int {
	return p.count
}
