package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawBamboolance struct {
	imgSizeX int
}

func (p *DrawBamboolance) Init() {
	var sx, sy int
	dxlib.GetGraphSize(imgBambooLance[0], &sx, &sy)

	p.imgSizeX = sx
}

func (p *DrawBamboolance) Draw(count int) {
	// Initを先に呼ばないと動かないようにする
	if p.imgSizeX <= 0 {
		common.SetError("実装にバグがあります。DrawBamboolance#Initを先に呼んでください")
		return
	}

	xreverse := int32(dxlib.TRUE)
	opt := dxlib.DrawRotaGraphOption{
		ReverseXFlag: &xreverse,
	}

	xd := count * 25
	if xd > battlecommon.PanelSize.X {
		xd = battlecommon.PanelSize.X
	}
	x := common.ScreenSize.X + p.imgSizeX/2 - xd
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		v := battlecommon.ViewPos(common.Point{X: 0, Y: y})
		dxlib.DrawRotaGraph(x, v.Y+battlecommon.PanelSize.Y/2, 1, 0, imgBambooLance[0], true, opt)
	}

}
