package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type Crack struct {
	Arg skillcore.Argument

	count     int
	attackPos []point.Point
}

func (p *Crack) Init(skillID int) {
	pos := p.Arg.GetObjectPos(p.Arg.OwnerID)

	switch skillID {
	case resources.SkillCrackout:
		p.attackPos = append(p.attackPos, point.Point{X: pos.X + 1, Y: pos.Y})
	case resources.SkillDoubleCrack:
		p.attackPos = append(p.attackPos, point.Point{X: pos.X + 1, Y: pos.Y})
		p.attackPos = append(p.attackPos, point.Point{X: pos.X + 2, Y: pos.Y})
	case resources.SkillTripleCrack:
		p.attackPos = append(p.attackPos, point.Point{X: pos.X + 1, Y: pos.Y - 1})
		p.attackPos = append(p.attackPos, point.Point{X: pos.X + 1, Y: pos.Y})
		p.attackPos = append(p.attackPos, point.Point{X: pos.X + 1, Y: pos.Y + 1})
	}
}

func (p *Crack) Process() (bool, error) {
	p.count++

	if p.count > 5 {
		p.Arg.SoundOn(resources.SEPanelBreak)
		for _, pos := range p.attackPos {
			p.Arg.PanelChange(pos, battlecommon.PanelStatusHole)
		}

		return true, nil
	}

	return false, nil
}

func (p *Crack) GetCount() int {
	return p.count
}
