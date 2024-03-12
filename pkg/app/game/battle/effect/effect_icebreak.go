package effect

import (
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
	/*
			その場: failed or 小爆発
		  右上: offY -30, ice_break[0]
		  右下: ice_break[1]
	*/
	return false, nil
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
}

func (e *iceBreakEffect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.ID,
		DrawType: anim.DrawTypeEffect,
	}
}
