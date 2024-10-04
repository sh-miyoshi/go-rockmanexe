package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type ChipForteAnother struct {
	Arg skillcore.Argument

	count    int
	state    int
	attacks  [4]ForteHellsRolling
	attackNo int
}

func (p *ChipForteAnother) Init() {
	p.setState(resources.SkillChipForteAnotherStateInit)
	for i := 0; i < 4; i++ {
		p.attacks[i] = ForteHellsRolling{Arg: p.Arg}
		if i%2 == 0 {
			p.attacks[i].Init(resources.SkillForteHellsRollingUp, true)
		} else {
			p.attacks[i].Init(resources.SkillForteHellsRollingDown, true)
		}
	}
}

func (p *ChipForteAnother) Process() (bool, error) {
	switch p.state {
	case resources.SkillChipForteAnotherStateInit:
		p.Arg.Cutin("フォルテアナザー", 300)
		p.setState(resources.SkillChipForteAnotherStateAppear)
		return false, nil
	case resources.SkillChipForteAnotherStateAppear:
		if p.count == 70 {
			p.setState(resources.SkillChipForteAnotherStateAttack)
			return false, nil
		}
	case resources.SkillChipForteAnotherStateAttack:
		return p.attacks[p.attackNo].Process()
	}
	p.count++
	return false, nil
}

func (p *ChipForteAnother) GetCount() int {
	return p.count
}

func (p *ChipForteAnother) GetState() int {
	return p.state
}

func (p *ChipForteAnother) GetAttackCount() int {
	return p.attacks[p.attackNo].GetCount()
}

func (p *ChipForteAnother) GetAttackPos() (prev point.Point, current point.Point, next point.Point) {
	return p.attacks[p.attackNo].GetPos()
}

func (p *ChipForteAnother) GetAttackNextStepCount() int {
	return forteHellsRollingNextStepCount
}

func (p *ChipForteAnother) setState(next int) {
	logger.Debug("Set ChipForteAnother state from %d to %d", p.state, next)
	p.count = 0
	p.state = next
}
