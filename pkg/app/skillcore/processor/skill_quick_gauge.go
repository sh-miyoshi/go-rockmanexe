package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type QuickGauge struct {
	Arg skillcore.Argument

	count int
}

func (p *QuickGauge) Update() (bool, error) {
	if p.count == 0 {
		p.Arg.Cutin("クイックゲージ", 90)
		battlecommon.CustomGaugeSpeed = 6
	}

	p.count++
	return p.count >= 90, nil
}

func (p *QuickGauge) GetCount() int {
	return p.count
}
