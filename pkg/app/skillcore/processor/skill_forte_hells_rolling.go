package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteActTypeUp = iota
	forteActTypeDown
)

const (
	nextStepCount = 20 // num(5) x delay(4)
)

type ForteHellsRolling struct {
	Arg skillcore.Argument

	count       int
	actType     int
	isFirstMove bool
	pos         point.Point
}

func (p *ForteHellsRolling) Init(skillID int) {
	p.count = 0
	p.isFirstMove = true
	p.pos = p.Arg.GetObjectPos(p.Arg.OwnerID)
	if skillID == resources.SkillForteHellsRollingUp {
		p.actType = forteActTypeUp
	} else {
		p.actType = forteActTypeDown
	}
}

func (p *ForteHellsRolling) Process() (bool, error) {
	p.count++
	if p.count%nextStepCount == 0 {
		if p.isFirstMove {
			switch p.actType {
			case forteActTypeUp:
				p.pos.Y--
			case forteActTypeDown:
				p.pos.Y++
			}
			p.isFirstMove = false
		}
		p.pos.X--
		if p.pos.X < 0 {
			return true, nil
		}

		// WIP: 一度だけプレイヤー方向に曲がる
		// WIP: ダメージ判定
	}
	return false, nil
}

func (p *ForteHellsRolling) GetCount() int {
	return p.count
}

func (p *ForteHellsRolling) GetPos() point.Point {
	return p.pos
}
