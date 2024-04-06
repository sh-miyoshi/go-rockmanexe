package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type AquamanWaterPipe struct {
	Arg skillcore.Argument
	Pos point.Point
}

type Aquaman struct {
	Arg skillcore.Argument

	count int
	state int
	pos   point.Point
	pipe  *AquamanWaterPipe
}

func (p *Aquaman) Init() {
	p.setState(resources.SkillAquamanStateInit)
	p.pos = p.Arg.GetObjectPos(p.Arg.OwnerID)
}

func (p *Aquaman) Process() (bool, error) {
	switch p.state {
	case resources.SkillAquamanStateInit:
		p.Arg.Cutin("アクアマン", 300)
		p.setState(resources.SkillAquamanStateAppear)
		return false, nil
	case resources.SkillAquamanStateAppear:
		if p.count == 70 {
			p.setState(resources.SkillAquamanStateCreatePipe)
			return false, nil
		}
	case resources.SkillAquamanStateCreatePipe:
		if p.count == 10 {
			p.pipe = &AquamanWaterPipe{
				Arg: p.Arg,
				Pos: point.Point{X: p.pos.X + 1, Y: p.pos.Y},
			}
			p.setState(resources.SkillAquamanStateAttack)
			return false, nil
		}
	case resources.SkillAquamanStateAttack:
		// 処理は呼び出し元で行う
		return false, nil
	}

	p.count++
	return false, nil
}

func (p *Aquaman) GetCount() int {
	return p.count
}

func (p *Aquaman) GetState() int {
	return p.state
}

func (p *Aquaman) GetPos() point.Point {
	return p.pos
}

func (p *Aquaman) PopWaterPipe() *AquamanWaterPipe {
	if p.pipe != nil {
		res := &AquamanWaterPipe{}
		*res = *p.pipe
		p.pipe = nil
		return res
	}
	return nil
}

func (p *Aquaman) setState(next int) {
	p.count = 0
	p.state = next
}
