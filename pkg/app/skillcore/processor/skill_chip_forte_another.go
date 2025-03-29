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

func (p *ChipForteAnother) Update() (bool, error) {
	switch p.state {
	case resources.SkillChipForteAnotherStateInit:
		if p.count == 0 {
			p.Arg.Cutin("フォルテアナザー", 500)
		}
		if p.count == 30 {
			p.Arg.MakeInvisible(p.Arg.OwnerID, 5) // ロックマンを消す
			p.setState(resources.SkillChipForteAnotherStateAppear)
			return false, nil
		}
	case resources.SkillChipForteAnotherStateAppear:
		if p.count == 90 {
			p.setState(resources.SkillChipForteAnotherStateAttack)
			return false, nil
		}
	case resources.SkillChipForteAnotherStateAttack:
		end, err := p.attacks[p.attackNo].Update()
		if err != nil {
			return false, err
		}
		if end {
			p.attackNo++
			if p.attackNo >= 4 {
				p.setState(resources.SkillChipForteAnotherStateEnd)
				return false, nil
			}
		}
	case resources.SkillChipForteAnotherStateEnd:
		if p.count == 30 {
			p.Arg.MakeInvisible(p.Arg.OwnerID, 0)
			return true, nil
		}
	}
	p.count++
	return false, nil
}

func (p *ChipForteAnother) GetCount() int {
	return p.count
}

func (p *ChipForteAnother) GetPos() point.Point {
	return p.Arg.GetObjectPos(p.Arg.OwnerID)
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

func (p *ChipForteAnother) setState(next int) {
	logger.Debug("Set ChipForteAnother state from %d to %d", p.state, next)
	p.count = 0
	p.state = next
}
