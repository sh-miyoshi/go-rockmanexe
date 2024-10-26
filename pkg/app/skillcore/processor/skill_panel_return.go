package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type PanelReturn struct {
	Arg   skillcore.Argument
	count int
}

func (p *PanelReturn) Process() (bool, error) {
	if p.count == 0 {
		p.Arg.Cutin("パネルリターン", 500)
	}

	if p.count == 50 {
		p.Arg.SoundOn(resources.SEPanelReturn)
	}

	p.count++
	return false, nil
}

func (p *PanelReturn) GetCount() int {
	return p.count
}
