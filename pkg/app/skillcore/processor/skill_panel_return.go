package processor

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"

type PanelReturn struct {
	Arg skillcore.Argument
}

func (p *PanelReturn) Process() (bool, error) {
	// WIP
	return true, nil
}

func (p *PanelReturn) GetCount() int {
	return 0
}
