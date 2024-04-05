package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type DrawBamboolance struct {
	imgSizeX int
}

func (p *DrawBamboolance) Init() {
	var sx, sy int
	dxlib.GetGraphSize(imgBambooLance[0], &sx, &sy)

	p.imgSizeX = sx
}

func (p *DrawBamboolance) Draw(count int, isPlayer bool) {
	// Initを先に呼ばないと動かないようにする
	if p.imgSizeX <= 0 {
		system.SetError("実装にバグがあります。DrawBamboolance#Initを先に呼んでください")
		return
	}

	// デフォルトが反転状態
	opt := dxlib.OptXReverse(isPlayer)
	xd := count * 25
	if xd > battlecommon.PanelSize.X {
		xd = battlecommon.PanelSize.X
	}
	x := config.ScreenSize.X + p.imgSizeX/2 - xd
	if !isPlayer {
		x = xd - p.imgSizeX/2
	}

	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		v := battlecommon.ViewPos(point.Point{X: 0, Y: y})
		dxlib.DrawRotaGraph(x, v.Y+battlecommon.PanelSize.Y/2, 1, 0, imgBambooLance[0], true, opt)
	}

}
