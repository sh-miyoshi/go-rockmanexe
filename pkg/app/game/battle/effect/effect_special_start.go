package effect

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	specialStartEffectEndCount = 20
)

type specialStartEffectLine struct {
	Direct     int
	StartPos   point.Point
	StartCount int
}

type specialStartEffect struct {
	ID  string
	Pos point.Point
	Ofs point.Point

	count int
	lines [4]specialStartEffectLine
}

func (p *specialStartEffect) Init() {
	view := battlecommon.ViewPos(p.Pos)
	p.lines[0] = specialStartEffectLine{StartPos: view, StartCount: 0, Direct: config.DirectRight | config.DirectDown}
	p.lines[1] = specialStartEffectLine{StartPos: view, StartCount: 5, Direct: config.DirectLeft | config.DirectDown}
	p.lines[2] = specialStartEffectLine{StartPos: view, StartCount: 10, Direct: config.DirectLeft | config.DirectUp}
	p.lines[3] = specialStartEffectLine{StartPos: view, StartCount: 15, Direct: config.DirectRight | config.DirectUp}
}

func (p *specialStartEffect) Process() (bool, error) {
	p.count++
	return p.count >= specialStartEffectEndCount, nil
}

func (p *specialStartEffect) Draw() {
	// 中央から拡大していく円
	view := battlecommon.ViewPos(p.Pos)
	x := view.X + p.Ofs.X
	y := view.Y + p.Ofs.Y
	r := p.count * 10
	const maxR = 30
	if r > maxR {
		r = maxR
	}
	dxlib.DrawCircle(x, y, r, dxlib.GetColor(255, 255, 255), false)

	// 4方向の直線
	for i := 0; i < 4; i++ {
		p.lines[i].Draw(p.count, 10)
	}
}

func (p *specialStartEffect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *specialStartEffectLine) Draw(count int, endCount int) {
	if count < p.StartCount {
		return
	}

	count = count - p.StartCount
	dx := 1
	dy := 1
	if p.Direct&config.DirectUp != 0 {
		dy = -1
	}
	if p.Direct&config.DirectLeft != 0 {
		dx = -1
	}

	switch count * 4 / endCount {
	case 0:
		dxlib.DrawLine(p.StartPos.X, p.StartPos.Y, p.StartPos.X+dx*10, p.StartPos.Y+dy*10, dxlib.GetColor(255, 255, 255))
	case 1:
		dxlib.DrawLine(p.StartPos.X, p.StartPos.Y, p.StartPos.X+dx*30, p.StartPos.Y+dy*30, dxlib.GetColor(255, 255, 255))
	case 2:
		dxlib.DrawLine(p.StartPos.X, p.StartPos.Y, p.StartPos.X+dx*80, p.StartPos.Y+dy*80, dxlib.GetColor(255, 255, 255))
	case 3:
		dxlib.DrawLine(p.StartPos.X+dx*50, p.StartPos.Y+dy*50, p.StartPos.X+dx*80, p.StartPos.Y+dy*80, dxlib.GetColor(255, 255, 255))
	}
}
