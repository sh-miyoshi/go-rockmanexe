package processor

import (
	"math/rand"

	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DeathMatch struct {
	SkillID int
	Arg     skillcore.Argument

	count     int
	breakList []point.Point
}

func (p *DeathMatch) Init() {
	p.count = 0
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		for x := 0; x < battlecommon.FieldNum.X; x++ {
			p.breakList = append(p.breakList, point.Point{X: x, Y: y})
		}
	}
	rand.Shuffle(len(p.breakList), func(i, j int) {
		p.breakList[i], p.breakList[j] = p.breakList[j], p.breakList[i]
	})
}

func (p *DeathMatch) Update() (bool, error) {
	if p.count == 0 {
		p.Arg.Cutin("デスマッチ", 500)
	}

	if p.count > 60 {
		if p.count%3 == 1 {
			if len(p.breakList) == 0 {
				return true, nil
			}

			pos := p.breakList[0]
			p.breakList = p.breakList[1:]

			var crackType int
			stEndCount := 0
			switch p.SkillID {
			case resources.SkillDeathMatch1:
				crackType = battlecommon.PanelStatusCrack
			case resources.SkillDeathMatch2:
				crackType = battlecommon.PanelStatusHole
				stEndCount = battlecommon.DefaultPanelStatusEndCount
			case resources.SkillDeathMatch3:
				crackType = battlecommon.PanelStatusPoison
			}

			p.Arg.ChangePanelStatus(pos, crackType, stEndCount)
			if p.count%9 == 1 {
				if p.SkillID == resources.SkillDeathMatch3 {
					p.Arg.SoundOn(resources.SEMakePoison)
				} else {
					p.Arg.SoundOn(resources.SEPanelBreakShort)
				}
			}
		}
	}

	p.count++
	return false, nil
}

func (p *DeathMatch) GetCount() int {
	return p.count
}
