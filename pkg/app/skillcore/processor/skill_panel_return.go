package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type PanelReturn struct {
	Arg     skillcore.Argument
	count   int
	targets []point.Point
}

func (p *PanelReturn) Update() (bool, error) {
	if p.count == 0 {
		p.targets = []point.Point{}
		p.Arg.Cutin("パネルリターン", 500)
	}

	if p.count == 50 {
		p.Arg.SoundOn(resources.SEPanelReturn)
		for x := 0; x < battlecommon.FieldNum.X; x++ {
			for y := 0; y < battlecommon.FieldNum.Y; y++ {
				pos := point.Point{X: x, Y: y}
				if p.Arg.GetPanelInfo(pos).Type == battlecommon.PanelTypePlayer {
					p.targets = append(p.targets, pos)
				}
			}
		}
		logger.Debug("Panel Return Target: %v", p.targets)
	}

	if p.count == 160 {
		for _, pos := range p.targets {
			p.Arg.ChangePanelStatus(pos, battlecommon.PanelStatusNormal, 0)
		}
		return true, nil
	}

	p.count++
	return false, nil
}

func (p *PanelReturn) GetCount() int {
	return p.count
}
