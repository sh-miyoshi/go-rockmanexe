package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	countBombEndCount = 30
)

type CountBomb struct {
	Arg skillcore.Argument

	count int
	pos   point.Point
}

func (p *CountBomb) Process() (bool, error) {
	if p.count == 0 {
		pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
		p.pos = point.Point{X: pos.X + 1, Y: pos.Y}
		p.Arg.Cutin("カウントボム")
	}

	p.count++
	return p.count > countBombEndCount, nil
}

func (p *CountBomb) GetCount() int {
	return p.count
}

func (p *CountBomb) GetEndCount() int {
	// EndCount = delay * (len(img) + keepCount)
	// TODO
	return 3 * 4
}

func (p *CountBomb) GetPos() point.Point {
	return p.pos
}
