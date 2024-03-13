package effect

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type iceBreakEffect struct {
	ID  string
	Pos point.Point
	Ofs point.Point

	count int
}

func (e *iceBreakEffect) Process() (bool, error) {
	e.count++
	return e.count >= len(images[resources.EffectTypeExplodeSmall]), nil
}

func (e *iceBreakEffect) Draw() {
	view := battlecommon.ViewPos(e.Pos)

	// 中央の小爆発
	imgNo := -1
	imgs := images[resources.EffectTypeExplodeSmall]
	delay := 1
	ofs := point.Point{}
	if e.count < len(imgs)*delay {
		imgNo = e.count / delay
	}

	if imgNo >= 0 {
		dxlib.DrawRotaGraph(view.X+ofs.X, view.Y+ofs.Y+15, 1, 0, imgs[imgNo], true)
	}

	// 右上へ流れていく破片
	delay = 40
	ofs = point.Point{X: e.count * delay, Y: -e.count*delay*2 - 20}
	x := view.X + ofs.X
	y := view.Y + ofs.Y
	const imgSizeX = 50
	const imgSizeY = 56
	if x < config.MaxScreenSize.X+imgSizeX/2 && y >= -imgSizeY/2 {
		dxlib.DrawRotaGraph(x, y, 1, 0, images[resources.EffectTypeIceBreak][0], true)
	}

	// 左下へ流れていく破片
	delay = 20
	ofs = point.Point{X: e.count * delay, Y: e.count*delay + 10}
	x = view.X + ofs.X
	y = view.Y + ofs.Y
	if x < config.MaxScreenSize.X+imgSizeX/2 && y < config.MaxScreenSize.Y+imgSizeY/2 {
		dxlib.DrawRotaGraph(x, y, 1, 0, images[resources.EffectTypeIceBreak][1], true)
	}
}

func (e *iceBreakEffect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.ID,
		DrawType: anim.DrawTypeEffect,
	}
}
