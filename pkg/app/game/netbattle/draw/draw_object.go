package draw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	appdraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/newnet/object"
)

func drawObject(images [object.TypeMax][]int, obj object.Object, opt Option) {
	view := battlecommon.ViewPos(common.Point{X: obj.X, Y: obj.Y})
	imgNo := obj.Count / object.ImageDelays[obj.Type]
	dxopts := dxlib.DrawRotaGraphOption{}

	if opt.Reverse {
		flag := int32(dxlib.TRUE)
		dxopts.ReverseXFlag = &flag
		obj.ViewOfsX *= -1
	}

	view.X += obj.ViewOfsX
	view.Y += obj.ViewOfsY

	// Special object draw
	switch obj.Type {
	case object.TypeVulcan:
		objectVulcan(images[obj.Type], view, imgNo, dxopts)
	case object.TypeWideShotMove:
		objectWideShotMove(images[obj.Type], view, obj, dxopts)
	case object.TypeThunderBall:
		objectThunderBall(images[obj.Type], view, obj, dxopts)
	case object.TypeMiniBomb:
		objectMiniBomb(images[obj.Type], view, obj, dxopts)
	default:
		if obj.Invincible {
			if cnt := obj.Count / 5 % 2; cnt == 0 {
				return
			}
		}

		if imgNo >= len(images[obj.Type]) {
			imgNo = len(images[obj.Type]) - 1
		}
		dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, images[obj.Type][imgNo], true, dxopts)
	}

	// Show HP
	if opt.ViewHP > 0 {
		appdraw.Number(view.X, view.Y+40, opt.ViewHP, appdraw.NumberOption{
			Color:    appdraw.NumberColorWhiteSmall,
			Centered: true,
		})
	}

	if len(obj.Chips) > 0 && opt.ViewChip {
		x := field.PanelSize.X*obj.X + field.PanelSize.X/2 - 18
		y := field.DrawPanelTopY + field.PanelSize.Y*obj.Y - 83
		dxlib.DrawBox(x-1, y-1, x+29, y+29, 0x000000, false)
		dxlib.DrawGraph(x, y, opt.ImgUnknownIcon, true)
	}
}

func objectVulcan(images []int, viewPos common.Point, imgNo int, dxopts dxlib.DrawRotaGraphOption) {
	if imgNo > 2 {
		imgNo /= 5 // slow down animation
	}

	ofsBody := 50
	ofsAtk := 100
	if dxopts.ReverseXFlag != nil && *dxopts.ReverseXFlag == dxlib.TRUE {
		ofsBody *= -1
		ofsAtk *= -1
	}

	// Show body
	no := imgNo
	if no > 2 {
		no = no % 2
	}

	dxlib.DrawRotaGraph(viewPos.X+ofsBody, viewPos.Y-18, 1, 0, images[no], true, dxopts)
	// Show attack
	if imgNo != 0 {
		if imgNo%2 == 0 {
			dxlib.DrawRotaGraph(viewPos.X+ofsAtk, viewPos.Y-10, 1, 0, images[3], true, dxopts)
		} else {
			dxlib.DrawRotaGraph(viewPos.X+ofsAtk, viewPos.Y-15, 1, 0, images[3], true, dxopts)
		}
	}
}

func objectWideShotMove(images []int, viewPos common.Point, obj object.Object, dxopts dxlib.DrawRotaGraphOption) {
	if obj.Speed == 0 {
		panic("ワイドショット描画のためのSpeedが0です")
	}

	imgNo := (obj.Count / object.ImageDelays[obj.Type]) % len(images)
	ofsx := field.PanelSize.X * obj.Count / obj.Speed
	if dxopts.ReverseXFlag != nil && *dxopts.ReverseXFlag == dxlib.TRUE {
		ofsx *= -1
	}

	dxlib.DrawRotaGraph(viewPos.X+ofsx, viewPos.Y, 1, 0, images[imgNo], true, dxopts)
}

func objectThunderBall(images []int, viewPos common.Point, obj object.Object, dxopts dxlib.DrawRotaGraphOption) {
	imgNo := (obj.Count / object.ImageDelays[obj.Type]) % len(images)

	if obj.Count >= obj.Speed {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	cnt := obj.Count % obj.Speed
	ofsx := battlecommon.GetOffset(obj.TargetX, obj.X, obj.PrevX, cnt, obj.Speed, field.PanelSize.X)
	ofsy := battlecommon.GetOffset(obj.TargetY, obj.Y, obj.PrevY, cnt, obj.Speed, field.PanelSize.Y)

	dxlib.DrawRotaGraph(viewPos.X+ofsx, viewPos.Y+25+ofsy, 1, 0, images[imgNo], true)
}

func objectMiniBomb(images []int, viewPos common.Point, obj object.Object, dxopts dxlib.DrawRotaGraphOption) {
	imgNo := (obj.Count / object.ImageDelays[obj.Type]) % len(images)

	// y = ax^2 + bx + c
	// (0,0), (d/2, ymax), (d, 0)
	size := field.PanelSize.X * 3
	ofsx := size * obj.Count / obj.Speed
	const ymax = 100
	ofsy := ymax*4*ofsx*ofsx/(size*size) - ymax*4*ofsx/size

	if dxopts.ReverseXFlag != nil && *dxopts.ReverseXFlag == dxlib.TRUE {
		ofsx *= -1
	}

	dxlib.DrawRotaGraph(viewPos.X+ofsx, viewPos.Y+ofsy, 1, 0, images[imgNo], true, dxopts)
}
